using System;
using System.IO;
using System.Runtime.InteropServices;
using ChillPatcher.Core; // 假定主项目的接口命名空间

namespace ChillPatcher.Module.QQMusic.CS
{
    // 示例结构体，代表音乐元信息
    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    public struct QQMusicMeta
    {
        public int Id;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 256)]
        public string Name;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 256)]
        public string Artist;
    }

    public class QQMusicModule : IMusicModule, IMusicSourceProvider
    {
        // DllImport 示例，将 Go 导出方法桥入
        [DllImport("qqmusic_bridge.dll", EntryPoint = "SearchMusic", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
        private static extern int SearchMusic(string keyword, out IntPtr metaArr, out int count);

        public string Name => "QQMusic C# Module";
        public string Description => "ChillPatcher插件，桥接Go实现的qqmusic_bridge。";

        public void Initialize()
        {
            // 初始化逻辑
        }

        public void Dispose()
        {
            // 清理逻辑
        }

        public MusicMeta[] Search(string keyword)
        {
            int ret = SearchMusic(keyword, out IntPtr metaPtr, out int count);
            if (ret != 0 || count == 0) return Array.Empty<MusicMeta>();

            var results = new MusicMeta[count];
            int structSize = Marshal.SizeOf<QQMusicMeta>();
            for (int i = 0; i < count; i++)
            {
                var meta = Marshal.PtrToStructure<QQMusicMeta>(metaPtr + i * structSize);
                results[i] = new MusicMeta(meta.Id, meta.Name, meta.Artist);
            }
            // 假设 Go 侧负责分配并需提供释放
            // [DllImport("qqmusic_bridge.dll")] static extern void FreeSearchResult(IntPtr ptr, int count);
            // FreeSearchResult(metaPtr, count);

            return results;
        }
    }
}
