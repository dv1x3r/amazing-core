import os
import sys
import json
import base64
import shutil
import struct
import hashlib
import argparse
import subprocess
import mutagen
import UnityPy
from PIL import Image
from collections import Counter

SIGNATURES = [
    ("audio/mp3", b"\xff\xfb"),
    ("audio/ogg", b"\x4f\x67\x67\x53\x00"),
    ("image/png", b"\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00"),
    ("AssetBundle/UnityFS", b"\x55\x6e\x69\x74\x79\x46\x53\x00"),
    ("AssetBundle/UnityWeb", b"\x55\x6e\x69\x74\x79\x57\x65\x62\x00"),
    ("TreeNode/Announcement", b"Announcement"),
    ("TreeNode/Areas", b"Areas"),
    ("TreeNode/AvatarProperty", b"AvatarProperty"),
    ("TreeNode/BuildingCompletion", b"BuildingCompletion"),
    ("TreeNode/BuildingUI", b"BuildingUI"),
    ("TreeNode/DressAvatarSlots", b"DressAvatarSlots"),
    ("TreeNode/Fish", b"Fish"),
    ("TreeNode/Game", b"Game"),
    ("TreeNode/Item", b"Item"),
    ("TreeNode/LevelUp", b"LevelUp"),
    ("TreeNode/Mission", b"Mission"),
    ("TreeNode/NPCs", b"NPCs"),
    ("TreeNode/NPCAnimations", b"NPCAnimations"),
    ("TreeNode/NPCRelationships", b"NPCRelationships"),
    ("TreeNode/Property", b"Property"),
    ("TreeNode/cQuest", b"cQuest"),
    ("TreeNode/Quest", b"Quest"),
    ("TreeNode/Root", b"Root"),
    ("TreeNode/SpawnPoints", b"SpawnPoints"),
    ("TreeNode/UIWidget", b"UIWidget"),
    ("json", b"{"),
]

UTF8_BOM = b"\xef\xbb\xbf"
VISUAL_TYPES = {"MeshFilter", "MeshRenderer", "SkinnedMeshRenderer"}


class UnitySceneParser:
    def __init__(self, env):
        self.env = env
        self.object_dict = {obj.path_id: obj for obj in env.objects}
        self.transform_dict = {}
        self.children_dict = {}
        self.roots = []
        for obj in self.env.objects:
            if obj.type.name == "Transform":
                tr = obj.parse_as_object()
                self.transform_dict[obj.path_id] = tr
                if not tr.m_Father:
                    self.roots.append(obj.path_id)
                else:
                    parent_id = tr.m_Father.path_id
                    self.children_dict.setdefault(parent_id, []).append(obj.path_id)

    def parse(self):
        return [self._build_tree(root_id) for root_id in self.roots]

    def _build_tree(self, tid, visited=None):
        if visited is None:
            visited = set()
        if tid in visited:
            return None
        visited.add(tid)

        tr = self.transform_dict[tid]
        path_id = tr.m_GameObject.path_id
        game_object = tr.m_GameObject.read()
        node = {"id": path_id, "name": game_object.m_Name}
        if path := tr.assets_file.container.path_dict.get(path_id):
            node["file"] = os.path.basename(path)
        if components := self._get_components(game_object):
            node["components"] = components
            if self._is_nonzero_transform(tr):
                node["transform"] = {
                    "position": {
                        "x": tr.m_LocalPosition.x,
                        "y": tr.m_LocalPosition.y,
                        "z": tr.m_LocalPosition.z,
                    },
                    "rotation": {
                        "x": tr.m_LocalRotation.x,
                        "y": tr.m_LocalRotation.y,
                        "z": tr.m_LocalRotation.z,
                        "w": tr.m_LocalRotation.w,
                    },
                    "scale": {
                        "x": tr.m_LocalScale.x,
                        "y": tr.m_LocalScale.y,
                        "z": tr.m_LocalScale.z,
                    },
                }
        children = []
        for child_id in self.children_dict.get(tid, []):
            child = self._build_tree(child_id, visited)
            if child is not None:
                children.append(child)
        if children:
            node["children"] = children
        return node

    def _is_nonzero_transform(self, tr):
        return any(
            [
                tr.m_LocalPosition.x,
                tr.m_LocalPosition.y,
                tr.m_LocalPosition.z,
                tr.m_LocalRotation.x,
                tr.m_LocalRotation.y,
                tr.m_LocalRotation.z,
                tr.m_LocalRotation.w,
                tr.m_LocalScale.x,
                tr.m_LocalScale.y,
                tr.m_LocalScale.z,
            ]
        )

    def _has_visual_components(self, game_object):
        for _, comp in game_object.m_Component:
            comp_obj = self.object_dict.get(comp.m_PathID)
            if comp_obj and comp_obj.type.name in VISUAL_TYPES:
                return True
        return False

    def _get_components(self, game_object):
        components = []
        for _, comp in game_object.m_Component:
            comp_obj = self.object_dict.get(comp.m_PathID)
            if comp_obj.type.name == "GameObject" or comp_obj.type.name == "Transform":
                continue
            component = {"type": comp_obj.type.name}
            if comp_obj.type.name == "AudioSource":
                data = comp_obj.parse_as_object()
                component["clip"] = self._parse_audio_source(data)
            elif comp_obj.type.name == "MonoBehaviour":
                data = comp_obj.parse_as_object()
                component["script"] = self._parse_mono_behaviour(data)
            elif comp_obj.type.name == "MeshFilter":
                data = comp_obj.parse_as_object()
                component["mesh"] = self._parse_mesh(data)
            elif comp_obj.type.name == "MeshRenderer":
                data = comp_obj.parse_as_object()
                if materials := self._parse_materials(data):
                    component["materials"] = materials
            elif comp_obj.type.name == "SkinnedMeshRenderer":
                data = comp_obj.parse_as_object()
                component["mesh"] = self._parse_mesh(data)
                if materials := self._parse_materials(data):
                    component["materials"] = materials
            elif comp_obj.type.name == "Animation":
                data = comp_obj.parse_as_object()
                if animations := self._parse_animations(data):
                    component["animations"] = animations
            components.append(component)
        return components

    def _resolve_ref(self, path_id, assets_file):
        ref = {"id": path_id}
        if obj := self.object_dict.get(path_id):
            ref["name"] = obj.peek_name()
        if path := assets_file.container.path_dict.get(path_id):
            ref["file"] = os.path.basename(path)
        return ref

    def _parse_audio_source(self, data):
        return self._resolve_ref(data.m_audioClip.m_PathID, data.assets_file)

    def _parse_mono_behaviour(self, data):
        script = self._resolve_ref(data.m_Script.m_PathID, data.assets_file)
        if hasattr(data, "triggerEvent"):
            script["trigger_event"] = data.triggerEvent
        return script

    def _parse_mesh(self, data):
        return self._resolve_ref(data.m_Mesh.m_PathID, data.assets_file)

    def _parse_material_textures(self, data):
        textures = []
        for tex_prop, tex_env in data.m_SavedProperties.m_TexEnvs:
            if tex_prop.name == "_MainTex":
                textures.append(
                    self._resolve_ref(tex_env.m_Texture.m_PathID, data.assets_file)
                )
        return textures

    def _parse_materials(self, data):
        materials = []
        for mat in data.m_Materials:
            mat_id = mat.m_PathID
            material = {"id": mat_id}
            if mat_obj := self.object_dict.get(mat_id):
                if mat_obj.type.name != "Material":
                    continue
                mat_data = mat_obj.parse_as_object()
                material["name"] = mat_data.m_Name
                if hasattr(mat_data, "m_SavedProperties"):
                    material["textures"] = self._parse_material_textures(mat_data)
            if mat_path := data.assets_file.container.path_dict.get(mat_id):
                material["file"] = os.path.basename(mat_path)
            materials.append(material)
        return materials

    def _parse_animations(self, data):
        return [
            self._resolve_ref(anim.m_PathID, data.assets_file)
            for anim in data.m_Animations
        ]


def get_bundle_info(env):
    info = {}

    for attr in (
        "signature",
        "version",
        "version_player",
        "version_engine",
    ):
        value = getattr(env.file, attr, None)
        if value is not None:
            info[attr] = str(value)

    # Fallback: derive version_engine from version
    if "version_engine" not in info:
        version = info.get("version")
        if version and version.startswith("UnityVersion "):
            info["version_engine"] = version.split(" ", 1)[1]

    return info


def get_bundle_counts(env):
    return {
        "assets": len(env.assets),
        "objects": len(env.objects),
        "container": len(env.container),
        "types": dict(sorted(Counter(obj.type.name for obj in env.objects).items())),
    }


def get_bundle_assets(env):
    return [
        {
            "name": asset.name,
            "target_platform": asset.target_platform.name,
        }
        for asset in env.assets
    ]


def get_bundle_containers(env):
    return {v.m_PathID: k for k, v in env.container.items()}


def unpack_assets(env, file_path, mp3_convert, zip_assets):
    base_name = os.path.basename(file_path)
    base_dir = os.path.dirname(file_path)
    assets_dir = os.path.join(base_dir, f"{base_name}_assets")
    unpacked_assets = 0

    for obj in env.objects:
        if obj.type.name == "AudioClip":
            clip = obj.parse_as_object()
            for sample_name, sample_data in clip.samples.items():
                os.makedirs(os.path.join(assets_dir, "audio"), exist_ok=True)
                tmp = os.path.join(assets_dir, "audio", sample_name)
                with open(tmp, "wb") as f:
                    f.write(sample_data)
                if mp3_convert:
                    print("    converting to mp3:", sample_name, file=sys.stderr)
                    file_name = f"{obj.path_id}.{os.path.splitext(sample_name)[0]}.mp3"
                    out = os.path.join(assets_dir, "audio", file_name)
                    subprocess.run(
                        ["ffmpeg", "-y", "-i", tmp, out],
                        capture_output=True,
                    )
                    os.remove(tmp)
                unpacked_assets += 1

        elif obj.type.name == "Mesh":
            data = obj.parse_as_object()
            file_name = f"{obj.path_id}.{data.m_Name}.obj"
            out = os.path.join(assets_dir, "models", file_name)
            if wavefront := data.export():
                os.makedirs(os.path.join(assets_dir, "models"), exist_ok=True)
                with open(out, "wt", newline="") as f:
                    f.write(wavefront)
                unpacked_assets += 1
            else:
                print("    invalid mesh:", data.m_Name, file=sys.stderr)

        elif obj.type.name == "Texture2D":
            data = obj.parse_as_object()
            file_name = f"{obj.path_id}.{data.m_Name}.png"
            out = os.path.join(assets_dir, "images", file_name)
            os.makedirs(os.path.join(assets_dir, "images"), exist_ok=True)
            try:
                data.image.save(out)
                unpacked_assets += 1
            except IsADirectoryError:
                print("    no image data:", data.m_Name, file=sys.stderr)

    if zip_assets and unpacked_assets > 0:
        zip_path = os.path.join(base_dir, f"{base_name}.zip")
        print("    compressing to zip:", zip_path, file=sys.stderr)
        shutil.make_archive(os.path.splitext(zip_path)[0], "zip", assets_dir)
        shutil.rmtree(assets_dir)


def get_file_info(file_path):
    file_name = os.path.basename(file_path)
    info = {
        "name": file_name,
        "oid": cdnid_to_oid(file_name),
        "type": detect_file_type(file_path),
        "hash": calculate_file_hash(file_path),
        "size": os.path.getsize(file_path),
    }
    return info


def get_audio_info(file_path):
    duration = get_audio_duration(file_path)
    if duration is None:
        print("    unable to get audio duration:", file_path, file=sys.stderr)
        duration = 0

    return {
        "duration": format_duration(duration),
    }


def get_image_info(file_path):
    dimensions = get_image_dimensions(file_path)
    if dimensions is None:
        print("    unable to get image dimensions:", file_path, file=sys.stderr)
        dimensions = "0x0"

    return {
        "dimensions": dimensions,
    }


def get_audio_duration(file_path):
    try:
        audio = mutagen.File(file_path)
    except Exception:
        return None

    if audio is None or not hasattr(audio, "info"):
        return None
    duration = getattr(audio.info, "length", None)
    if duration is None:
        return None
    return duration


def get_image_dimensions(file_path):
    try:
        with Image.open(file_path) as image:
            width, height = image.size
    except Exception:
        return None

    return f"{width}x{height}"


def format_duration(seconds):
    total_seconds = max(0, int(round(seconds)))
    minutes, secs = divmod(total_seconds, 60)
    hours, minutes = divmod(minutes, 60)
    if hours:
        return f"{hours:02}:{minutes:02}:{secs:02}"
    return f"{minutes:02}:{secs:02}"


def detect_file_type(file_path):
    with open(file_path, "rb") as f:
        raw_sample = f.read(256)

    sample = raw_sample.strip().removeprefix(UTF8_BOM)

    for type_name, signature in SIGNATURES:
        if sample.startswith(signature):
            return type_name

    if is_unity_serialized_file(file_path, raw_sample):
        return "Unity/SerializedFile"

    return "Unknown"


def is_unity_serialized_file(file_path, sample):
    if len(sample) < 16:
        return False

    metadata_size, file_size, version, data_offset = struct.unpack(">IIII", sample[:16])
    actual_size = os.path.getsize(file_path)

    if version < 5 or version > 30:
        return False
    if file_size < actual_size:
        return False
    if metadata_size <= 0 or metadata_size >= actual_size:
        return False
    if version >= 9 and (data_offset < metadata_size or data_offset > actual_size):
        return False

    return True


def calculate_file_hash(file_path):
    with open(file_path, "rb") as f:
        return hashlib.sha1(f.read()).hexdigest()


def cdnid_to_oid(cdnid):
    try:
        padded = cdnid + "=" * (-len(cdnid) % 4)
        decoded = base64.b64decode(padded)
        return decoded.decode("ascii")
    except:
        print("    unable to get oid from cdnid:", cdnid, file=sys.stderr)
        return 0


def process_file(file_path, unpack=None, mp3_convert=None, zip_assets=None):
    file_info = get_file_info(file_path)
    metadata = {"file": file_info}
    if file_info["type"] in ("audio/mp3", "audio/ogg"):
        metadata["audio"] = get_audio_info(file_path)
    elif file_info["type"] == "image/png":
        metadata["image"] = get_image_info(file_path)
    elif file_info["type"] in [
        "AssetBundle/UnityFS",
        "AssetBundle/UnityWeb",
        "Unity/SerializedFile",
    ]:
        with open(file_path, "rb") as f:
            env = UnityPy.load(f)
            bundle = {}
            bundle["info"] = get_bundle_info(env)
            bundle["assets"] = get_bundle_assets(env)
            bundle["counts"] = get_bundle_counts(env)
            bundle["containers"] = get_bundle_containers(env)
            scene = UnitySceneParser(env)
            print("    parsing scene...", file=sys.stderr)
            bundle["scene"] = scene.parse()
            if unpack:
                print("    unpacking assets...", file=sys.stderr)
                unpack_assets(env, file_path, mp3_convert, zip_assets)
            metadata["bundle"] = bundle
    return metadata


def write_metadata(metadata, json_path):
    if json_dir := os.path.dirname(json_path):
        os.makedirs(json_dir, exist_ok=True)
    with open(json_path, "w") as f:
        json.dump(metadata, f, indent=2)


def should_skip_path(path):
    name = os.path.basename(path).lower()
    return (
        name.endswith(".meta.json") or name.endswith(".zip") or name.endswith("_assets")
    )


def process_manifest(manifest_path, unpack=None, mp3_convert=None, zip_assets=None):
    with open(manifest_path) as f:
        manifest = json.load(f)

    files = manifest.get("files", [])
    for i, entry in enumerate(files, 1):
        file_path = entry["file_path"]
        json_path = entry["json_path"]
        print(f"processing file {i}/{len(files)}: {file_path}", file=sys.stderr)

        try:
            metadata = process_file(
                file_path,
                unpack=unpack,
                mp3_convert=mp3_convert,
                zip_assets=zip_assets,
            )
            write_metadata(metadata, json_path)
        except Exception as err:
            raise RuntimeError(f"failed to process file: {file_path}") from err

        print("    metadata written:", json_path, file=sys.stderr)


def main():
    parser = argparse.ArgumentParser(
        prog="cache.py",
        description="Inspect and unpack game cache files.",
    )
    parser.add_argument(
        "path",
        nargs="?",
        help="Path to a single cache file or a folder containing multiple files.",
    )
    parser.add_argument(
        "--manifest",
        help="JSON manifest with file_path/json_path entries.",
    )
    parser.add_argument(
        "--stdout",
        action="store_true",
        help="Write JSON metadata to stdout.",
    )
    parser.add_argument(
        "--json",
        help="Write JSON metadata to .meta.json files.",
        action="store_true",
    )
    parser.add_argument(
        "--unpack",
        action="store_true",
        help="Unpack assets (audio, images, models).",
    )
    parser.add_argument(
        "--mp3",
        action="store_true",
        help="Convert audio to mp3 using ffmpeg (requires --unpack).",
    )
    parser.add_argument(
        "--zip",
        action="store_true",
        help="Zip unpacked assets (requires --unpack).",
    )

    args = parser.parse_args()

    if args.manifest:
        if args.path:
            parser.error("path cannot be used with --manifest")
        if args.stdout:
            parser.error("--stdout cannot be used with --manifest")
        if args.json:
            parser.error("--json cannot be used with --manifest")
        if (args.mp3 or args.zip) and not args.unpack:
            parser.error("--mp3 and --zip require --unpack")
        process_manifest(
            args.manifest,
            unpack=args.unpack,
            mp3_convert=args.mp3,
            zip_assets=args.zip,
        )
        return

    if not args.path:
        parser.error("path is required unless --manifest is used")

    if (args.mp3 or args.zip) and not args.unpack:
        parser.error("--mp3 and --zip require --unpack")

    if not args.stdout and not args.json and not args.unpack:
        parser.error("at least one of --stdout, --json or --unpack is required")

    if should_skip_path(args.path):
        files = []
    elif os.path.isdir(args.path):
        files = sorted(
            [
                os.path.join(args.path, f)
                for f in os.listdir(args.path)
                if os.path.isfile(os.path.join(args.path, f))
                and not should_skip_path(os.path.join(args.path, f))
            ]
        )
    else:
        files = [args.path]

    metadata_entries = []
    for i, file_path in enumerate(files, 1):
        print(f"processing file {i}/{len(files)}: {file_path}", file=sys.stderr)

        metadata = process_file(
            file_path,
            unpack=args.unpack,
            mp3_convert=args.mp3,
            zip_assets=args.zip,
        )
        metadata_entries.append(metadata)

        if args.json:
            out_path = file_path + ".meta.json"
            write_metadata(metadata, out_path)
            print("    metadata written:", out_path, file=sys.stderr)

    if args.stdout:
        value = metadata_entries[0] if len(metadata_entries) == 1 else metadata_entries
        print(json.dumps(value, separators=(",", ":")))


if __name__ == "__main__":
    main()
