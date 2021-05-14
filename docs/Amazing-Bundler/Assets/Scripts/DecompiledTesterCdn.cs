using System.Collections.Generic;
using UnityEngine;

public class DecompiledTesterCdn : MonoBehaviour
{
    // Start is called before the first frame update
    void Start()
    {
        var asset = new GSFAsset();
        asset.assetTypeName = "asset_type";
        asset.cdnId = "Player_Base.unity3d";
        asset.resName = "res_name";
        asset.groupName = "group_name";
        asset.fileSize = 1079; // size in bytes

        AssetDownloadManager.LoadAssetMap(new Dictionary<string, List<GSFAsset>>() {
            { "Preload_PrefabUnity3D", new List<GSFAsset>() { asset } }
        }, null);
        DownloadManager.Instance.Update();
    }

    // Update is called once per frame
    void Update()
    {

    }
}
