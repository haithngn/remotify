# Code Generation Templates

This directory contains template files for the Generate Code feature in Remotify.

## Structure

```
templates/
├── client/          # Client-side notification code templates
│   ├── ios-swift.tmpl
│   ├── ios-objc.tmpl
│   ├── android-kotlin.tmpl
│   ├── android-java.tmpl
│   ├── flutter.tmpl
│   └── react-native.tmpl
└── server/          # Server-side notification code templates
    ├── nodejs.tmpl
    ├── python.tmpl
    ├── go.tmpl
    ├── java.tmpl
    ├── csharp.tmpl
    └── php.tmpl
```

## Template Format

Templates use a simple placeholder system with double curly braces:

- `{{title}}` - Notification title
- `{{body}}` - Notification body

Example:

```swift
let content = UNMutableNotificationContent()
content.title = "{{title}}"
content.body = "{{body}}"
```

## Adding New Templates

1. Create a new `.tmpl` file in the appropriate directory (client/ or server/)
2. Use the `{{title}}` and `{{body}}` placeholders where appropriate
3. Add the new platform/language option to the dropdown in `GenerateCode.vue`
4. The template loader will automatically pick up the new template

## Template Loading

Templates are loaded dynamically using Vite's raw import feature:

- Templates are cached after first load for performance
- Loading is asynchronous with proper error handling
- Templates are processed with variable substitution when generating code

## Supported Platforms

### Client Platforms

- **iOS Swift** - Local notification using UserNotifications framework
- **iOS Objective-C** - Local notification using UserNotifications framework
- **Android Kotlin** - Local notification using NotificationCompat
- **Android Java** - Local notification using NotificationCompat
- **Flutter** - Local notification using flutter_local_notifications plugin
- **React Native** - Local notification using react-native-push-notification

### Server Languages

- **Node.js** - Firebase Admin SDK for JavaScript
- **Python** - Firebase Admin SDK for Python
- **Go** - Firebase Admin SDK for Go
- **Java** - Firebase Admin SDK for Java
- **C#** - Firebase Admin SDK for .NET
- **PHP** - Firebase Admin SDK for PHP

All server templates assume Firebase Cloud Messaging (FCM) for cross-platform push notifications.

## Production Build Notes

**Important**: Templates are loaded using a static import registry (`templateRegistry.ts`) for production builds:
- All templates are imported at build time for optimal performance
- Templates are bundled into the application (no runtime file loading)
- When adding new templates, remember to update `templateRegistry.ts`

### Template Registry Structure
```typescript
// Import templates
import newTemplate from '../templates/client/new-platform.tmpl?raw'

// Export in registry
export const clientTemplates: { [key: string]: string } = {
  'new-platform': newTemplate,
  // ... other templates
}
```
