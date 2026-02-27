# CorpFlow Build Scripts

## 一键构建 (需要Flutter环境)

```bash
# 进入移动端目录
cd mobile

# 安装依赖
flutter pub get

# 构建 Android APK
flutter build apk --debug

# 构建 Android Release APK
flutter build apk --release

# 构建 iOS (macOS only)
flutter build ios --release

# 构建 macOS
flutter build macos --release

# 构建 Windows
flutter build windows --release

# 构建 Web
flutter build web
```

## 国产Android应用商店发布

### 1. 华为应用市场
- 注册开发者账号
- 上传签名后的APK
- 填写应用信息

### 2. 小米应用商店
- 使用小米开发者账号
- 上传APK+icon+截图

### 3. OPPO/vivo应用商店
- 各自开发者平台上传

## 跨平台构建配置

```yaml
# pubspec.yaml 已配置:
name: corpflow
description: CorpFlow Mobile Client

# 支持的平台:
# - Android (APK/AAB)
# - iOS (IPA)
# - macOS (APP)
# - Windows (EXE)
# - Linux (AppImage/DEB)
# - Web (HTML5)
```

## 环境要求

```bash
# Flutter SDK
# https://flutter.dev/docs/get-started/install

# Android 开发
# - Android Studio
# - JDK 17+
# - Android SDK

# iOS 开发 (macOS only)
# - Xcode 15+
# - CocoaPods
```
