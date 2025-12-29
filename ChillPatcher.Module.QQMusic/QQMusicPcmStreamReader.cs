using System;
using ChillPatcher.SDK.Interfaces;
using ChillPatcher.SDK.Models;

namespace ChillPatcher.Module.QQMusic
{
    public class QQMusicPcmStreamReader : IPcmStreamReader
    {
        // [TODO] 仿照网易云实现所有接口
        public PcmStreamInfo Info => throw new NotImplementedException();
        public ulong CurrentFrame => throw new NotImplementedException();
        public bool IsEndOfStream => throw new NotImplementedException();
        public bool IsReady => throw new NotImplementedException();
        public bool WaitForReady(int timeoutMs = 5000) => throw new NotImplementedException();
        public int Read(float[] buffer, int offset, int count) => throw new NotImplementedException();
        public void Dispose() { }
    }
}
