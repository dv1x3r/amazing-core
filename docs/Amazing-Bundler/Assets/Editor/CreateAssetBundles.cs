using UnityEngine;
using UnityEditor;
using System.IO;
using System.Linq;

public class CreateAssetBundles
{
    [MenuItem("Assets/Build AssetBundles From Properties - StreamingAssets (no mainAsset)")]
    static void BuildAllAssetBundles()
    {
        string assetBundleDirectory = "Assets/StreamingAssets";
        if (!Directory.Exists(Application.streamingAssetsPath))
        {
            Directory.CreateDirectory(assetBundleDirectory);
        }
        BuildPipeline.BuildAssetBundles(assetBundleDirectory, BuildAssetBundleOptions.None, EditorUserBuildSettings.activeBuildTarget);
    }

    [MenuItem("Assets/Build AssetBundle From Selection - Track dependencies (with mainAsset)")]
    static void ExportResource()
    {
        // Bring up save panel
        string path = EditorUtility.SaveFilePanel("Save Resources", "", "New Resource", "unity3d");
        if (path.Length != 0)
        {
            // Build the resource file from the active selection.
            Object[] selection = Selection.GetFiltered(typeof(Object), SelectionMode.DeepAssets);
            Debug.Log("Exporting main asset: [" + Selection.activeObject.name + "]    Selection: [" +
                      string.Join("], [", (from item in selection select item.name).ToArray()) + "]");
            BuildPipeline.BuildAssetBundle(Selection.activeObject, selection, path,
                                           BuildAssetBundleOptions.CollectDependencies | BuildAssetBundleOptions.CompleteAssets,
                                           BuildTarget.StandaloneWindows);
            Selection.objects = selection;
        }
    }
}
