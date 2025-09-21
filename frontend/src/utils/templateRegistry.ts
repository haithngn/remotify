// Template registry - imports all templates as modules for production builds

// Client templates
import iosSwiftTemplate from "../templates/client/ios-swift.tmpl?raw";
import iosObjcTemplate from "../templates/client/ios-objc.tmpl?raw";
import androidKotlinTemplate from "../templates/client/android-kotlin.tmpl?raw";
import androidJavaTemplate from "../templates/client/android-java.tmpl?raw";
import flutterTemplate from "../templates/client/flutter.tmpl?raw";
import reactNativeTemplate from "../templates/client/react-native.tmpl?raw";

// Server templates
import nodejsTemplate from "../templates/server/nodejs.tmpl?raw";
import pythonTemplate from "../templates/server/python.tmpl?raw";
import goTemplate from "../templates/server/go.tmpl?raw";
import javaTemplate from "../templates/server/java.tmpl?raw";
import csharpTemplate from "../templates/server/csharp.tmpl?raw";
import phpTemplate from "../templates/server/php.tmpl?raw";
import rubyTemplate from "../templates/server/ruby.tmpl?raw";

export const clientTemplates: { [key: string]: string } = {
  "ios-swift": iosSwiftTemplate,
  "ios-objc": iosObjcTemplate,
  "android-kotlin": androidKotlinTemplate,
  "android-java": androidJavaTemplate,
  flutter: flutterTemplate,
  "react-native": reactNativeTemplate,
};

export const serverTemplates: { [key: string]: string } = {
  nodejs: nodejsTemplate,
  python: pythonTemplate,
  go: goTemplate,
  java: javaTemplate,
  csharp: csharpTemplate,
  php: phpTemplate,
  ruby: rubyTemplate
};
