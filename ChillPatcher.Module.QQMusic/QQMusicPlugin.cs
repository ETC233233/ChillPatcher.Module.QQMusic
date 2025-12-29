using BepInEx;
using System;
using System.IO;

namespace ChillPatcher.Module.QQMusic
{
    [BepInPlugin("chillpatcher.qqmusic", "ChillPatcher QQMusic Module", "0.1.0")]
    public class QQMusicPlugin : BaseUnityPlugin
    {
        private void Awake()
        {
            Logger.LogInfo("[QQMusic] 模块已加载！");
            QQMusicBridge.Initialize(Logger);
        }
    }
}
