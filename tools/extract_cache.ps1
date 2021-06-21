$ripper = "C:\Users\dx\source\repos\UtinyRipper-master\Bins\CP\Debug\uTinyRipperCmd.dll"
$src_dir = "C:\Users\dx\source\repos\UtinyRipper-master\Bins\CP\Debug\Cache\"
$target_dir = "C:\Users\dx\source\repos\UtinyRipper-master\Bins\CP\Debug\"
cd $target_dir

Get-ChildItem $src_dir | Foreach-Object {
    dotnet $ripper $_.FullName
}
