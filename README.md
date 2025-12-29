# ChillPatcher.Module.QQMusic

QQ音乐模块，扫码登录、歌单/音源等全部对标新浪潮网易云接口实现（重点参考 copws/qq-music-api ），可供 BepInEx 与 ChillPatcher 框架直接集成。

## 目录说明
- `qqmusic_bridge/` Go桥接，负责账号扫码、歌单、音源接口对接QQ音乐API —— 详见 copws/qq-music-api。
- `ChillPatcher.Module.QQMusic/` C#桥接与PCM等，接口仿照 ChillPatcher.Module.Netease。

## 开发要点/TODO
- [ ] Go桥接建议直接模仿 copws/qq-music-api 结构，从扫码/login/qr等接口入手
- [ ] C#部分只需调用Go导出的api，流式音频拉取与歌单同步实现见网易云模块
- [ ] 任何难以Go落地的部分，建议先用Node/微服务方案辅助
