#!/bin/bash

# CorpFlow Mobile Build Script
# 一键构建移动应用

set -e

echo "🚀 CorpFlow Mobile Build Script"
echo "================================"

# 检查Flutter
if ! command -v flutter &> /dev/null; then
    echo "❌ Flutter 未安装"
    echo "请先安装 Flutter: https://flutter.dev/docs/get-started/install"
    exit 1
fi

# 显示Flutter版本
echo "📱 Flutter version: $(flutter --version | head -1)"
echo ""

# 进入移动端目录
cd "$(dirname "$0")"

# 获取依赖
echo "📦 Installing dependencies..."
flutter pub get
echo ""

# 创建输出目录
mkdir -p output

# 构建函数
build_android() {
    echo "🔨 Building Android APK..."
    flutter build apk --release
    cp build/app/outputs/flutter-apk/app-release.apk output/corpflow-android.apk
    echo "✅ Android APK: output/corpflow-android.apk"
}

build_ios() {
    echo "🔨 Building iOS..."
    flutter build ios --release --no-codesign
    echo "✅ iOS build complete: build/ios/iphoneos/"
}

build_macos() {
    echo "🔨 Building macOS..."
    flutter build macos --release
    echo "✅ macOS build complete: build/macos/Build/Products/Release/"
}

build_windows() {
    echo "🔨 Building Windows..."
    flutter build windows --release
    cp build/windows/x64/runner/Release/CorpFlow.exe output/
    echo "✅ Windows EXE: output/CorpFlow.exe"
}

build_web() {
    echo "🔨 Building Web..."
    flutter build web --release
    echo "✅ Web build: build/web/"
}

# 根据参数选择构建目标
case "${1:-all}" in
    android)
        build_android
        ;;
    ios)
        build_ios
        ;;
    macos)
        build_macos
        ;;
    windows)
        build_windows
        ;;
    web)
        build_web
        ;;
    all)
        echo "Building all platforms..."
        echo ""
        build_android
        # 注意: iOS/macOS只能在macOS上构建
        # Windows只能在Windows上构建
        if [[ "$OSTYPE" == "darwin"* ]]; then
            build_ios
            build_macos
        elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
            build_windows
        fi
        build_web
        ;;
    *)
        echo "Usage: $0 [android|ios|macos|windows|web|all]"
        echo ""
        echo "Examples:"
        echo "  $0 android    # Build Android APK only"
        echo "  $0 all        # Build all platforms"
        exit 1
        ;;
esac

echo ""
echo "🎉 Build complete!"
