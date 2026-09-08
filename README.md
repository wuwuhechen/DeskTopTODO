# DesktopTODO

基于 Go、Wails v3 与 React 的桌面待办便签应用。

应用以无边框、淡紫液态玻璃卡片的形式显示在桌面上，支持本地待办管理、窗口状态保存、系统托盘与开机启动。

本项目已完成windows平台的开发与打包，可跳转到release进行下载与使用。

## 功能特性

- 待办事项新增、编辑、完成与删除
- 高、中、低、无优先级标记
- 已完成事项显示与隐藏
- 待办事项拖拽排序
- SQLite 本地数据持久化
- 窗口位置、尺寸、置顶和锁定状态保存
- 淡紫液态玻璃外观
- 系统托盘显示、隐藏与退出
- Windows 登录后静默启动至托盘

## 技术栈

- Go
- Wails v3 Beta
- React + TypeScript + Vite
- GORM + SQLite
- dnd-kit
- NSIS

## 项目结构

```text
DesktopTODO/
├── frontend/                     # React 前端
│   ├── src/                      # 页面、组件与样式
│   │   ├── App.tsx               # 页面状态与交互
│   │   ├── App.css               # 液态玻璃样式
│   │   ├── todo_item.tsx         # 单条待办组件
│   │   └── todo_list.tsx         # 拖拽排序列表
│   └── bindings/                 # Wails 自动生成的前端绑定
├── internal/
│   ├── database/                 # 数据库初始化
│   ├── model/                    # Todo、设置和窗口状态模型
│   └── repository/               # GORM 数据访问层
├── build/
│   ├── config.yml                # 产品元数据
│   └── windows/                  # 图标、NSIS 与 Windows 打包配置
├── main.go                       # Wails 应用与主窗口入口
├── todoservice.go                # 待办业务服务
├── settingservice.go             # 用户偏好服务
└── windowservice.go              # 窗口、托盘与开机启动服务
```

## 环境要求

- Windows 11（推荐，用于原生半透明与 Acrylic 背景效果）
- Go 1.25 或更高版本
- Node.js 与 npm
- [Wails v3 CLI](https://v3.wails.io/getting-started/installation/)
- WebView2 Runtime

安装 Wails CLI：

```powershell
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 doctor
```

## 快速开始

1. 克隆项目到本地

   ```bash
   git clone <repository-url>
   cd DesktopTODO
   ```

2. 安装前端依赖

   ```powershell
   cd frontend
   npm install
   cd ..
   ```

3. 启动开发模式

   ```powershell
   wails3 dev
   ```

4. 应用启动后，即可在桌面上新增、编辑、完成和排序待办事项。

## 数据与设置

待办、用户偏好和窗口状态均保存至用户本地 SQLite 数据库。应用不需要登录，也不会默认上传数据。

主要用户设置包括：

- 是否显示已完成事项
- 是否始终置顶
- 是否锁定窗口位置
- 是否开机启动
- 上一次窗口的位置和尺寸

## Windows 打包

1. 安装 NSIS

   ```powershell
   winget install -e --id NSIS.NSIS --source winget
   ```

2. 如果 PowerShell 找不到 `makensis`，将其加入当前终端路径

   ```powershell
   $env:Path += ";C:\Program Files (x86)\NSIS"
   makensis /VERSION
   ```

3. 打包当前用户可安装的版本

   ```powershell
   wails3 package INSTALL_SCOPE=user
   ```

4. 构建产物位于 `bin/` 目录

   ```text
   bin/desktoptodo.exe
   bin/desktoptodo-amd64-installer.exe
   ```

> [!warning]
> Wails v3 目前仍处于 Beta 阶段。Windows 安装程序使用简体中文界面，并显式以 UTF-8 编码编译 NSIS 脚本。

## Android 打包

仓库包含 Wails v3 Android 模板。Android 打包需要 Android SDK、NDK 26.3.x、JDK 与签名配置，建议在 WSL2/Linux 环境中执行：

```bash
wails3 task android:package  # 生成 APK
wails3 task android:bundle   # 生成 Play Store 所需 AAB
```

Android 不支持桌面系统托盘、窗口置顶和桌面玻璃窗口，发布前需要进行移动端适配。

## 后续计划

- 支持多便签
- 添加截止日期和本地提醒
- 实现数据导入、导出与备份
- 完善 Android 专属交互与发布配置
- 添加自动更新和代码签名
