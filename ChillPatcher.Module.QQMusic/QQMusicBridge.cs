using System;
using System.Runtime.InteropServices;
using BepInEx.Logging;

namespace ChillPatcher.Module.QQMusic
{
    public static class QQMusicBridge
    {
        private static ManualLogSource _logger;
        public static void Initialize(ManualLogSource logger)
        {
            _logger = logger;
            _logger.LogInfo("[QQMusicBridge] Initialized");
            // [TODO] 动态加载 qqmusic_bridge DLL
        }

        // 扫码登录相关API
        [DllImport("qqmusic_bridge", CallingConvention = CallingConvention.Cdecl)]
        public static extern IntPtr QQMusicQRGetKey();
        [DllImport("qqmusic_bridge", CallingConvention = CallingConvention.Cdecl)]
        public static extern IntPtr QQMusicQRCheckStatus();
        [DllImport("qqmusic_bridge", CallingConvention = CallingConvention.Cdecl)]
        public static extern void QQMusicQRCancelLogin();
        // [TODO] 更多API如歌单/音频拉取...
        public static string GetLoginQRCode()
        {
            IntPtr ptr = QQMusicQRGetKey();
            var json = Marshal.PtrToStringAnsi(ptr);
            return json;
        }
    }
}
