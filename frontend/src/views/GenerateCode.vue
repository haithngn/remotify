<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import { GetSettings } from "../../wailsjs/go/main/App";
import { values } from "../../wailsjs/go/models";
import { templateLoader } from "../utils/templateLoader";
import { notify } from "../components/notification";

// Define props to receive payload data from parent tabs
const props = defineProps({
  apnsPayload: {
    type: String,
    default: "",
  },
  fcmPayload: {
    type: String,
    default: "",
  },
});

// Client Platform options
const clientPlatforms = [
  { value: "ios-swift", label: "iOS Swift" },
  { value: "ios-objc", label: "iOS Objective-C" },
  { value: "android-kotlin", label: "Android Kotlin" },
  { value: "android-java", label: "Android Java" },
  { value: "flutter", label: "Flutter" },
  { value: "react-native", label: "React Native" },
];

// Server Language options
const serverLanguages = [
  { value: "nodejs", label: "Node.js" },
  { value: "python", label: "Python" },
  { value: "go", label: "Go" },
  { value: "java", label: "Java" },
  { value: "csharp", label: "C#" },
  { value: "php", label: "PHP" },
  { value: "ruby", label: "Ruby" },
];

// Reactive state
const defaultPayload = {
  aps: {
    alert: {
      title: "Notification Title",
      subtitle: "Notification Subtitle",
      body: "Notification Message Body",
    },
    sound: "default",
    badge: 19,
    category: "CATEGORY_IDENTIFIER_SHOULD_BE_UPPERCASED",
    "thread-id": "message_thread_id",
    "mutable-content": 1,
  },
  custom_data: {
    key: "value",
    nested_object: {
      key: "value",
    },
  },
};
const selectedClientPlatform = ref("ios-swift");
const selectedServerLanguage = ref("nodejs");
const editPayloadHere = ref(false);
const customPayload = ref<any>(defaultPayload);
const customPayloadText = ref(JSON.stringify(customPayload.value, null, 2));
const clientCodeTemplate = ref("");
const serverCodeTemplate = ref("");
const isLoadingClientTemplate = ref(false);
const isLoadingServerTemplate = ref(false);

const isDarkMode = ref(false);

// Load theme on mount
onMounted(() => {
  GetSettings().then((data) => {
    let settings = JSON.parse(data);
    isDarkMode.value = settings.data.theme_mode === "dark";
  });

  // Load initial templates
  loadClientTemplate();
  loadServerTemplate();
});

// Load client template
const loadClientTemplate = async () => {
  isLoadingClientTemplate.value = true;
  try {
    const template = await templateLoader.loadClientTemplate(
      selectedClientPlatform.value
    );
    clientCodeTemplate.value = template;
  } catch (error) {
    console.error("Failed to load client template:", error);
    clientCodeTemplate.value = "// Failed to load template";
  } finally {
    isLoadingClientTemplate.value = false;
  }
};

// Load server template
const loadServerTemplate = async () => {
  isLoadingServerTemplate.value = true;
  try {
    const template = await templateLoader.loadServerTemplate(
      selectedServerLanguage.value
    );
    serverCodeTemplate.value = template;
  } catch (error) {
    console.error("Failed to load server template:", error);
    serverCodeTemplate.value = "// Failed to load template";
  } finally {
    isLoadingServerTemplate.value = false;
  }
};

// Watch for platform/language changes
watch(selectedClientPlatform, loadClientTemplate);
watch(selectedServerLanguage, loadServerTemplate);

// Generated client code
const clientCode = computed(() => {
  if (isLoadingClientTemplate.value) {
    return "// Loading template...";
  }

  if (!clientCodeTemplate.value) {
    return "// Template not available";
  }

  let payload: any = customPayload.value;

  // If not editing payload here, use the default template or parent payload
  if (!editPayloadHere.value) {
    try {
      // Use APNS payload if available, otherwise use default
      if (props.apnsPayload) {
        payload = JSON.parse(props.apnsPayload);
      } else {
        payload = JSON.parse(values.PayloadTemplate.DefaultAPNS);
      }
    } catch (e) {
      // Fallback to current payload if parsing fails
      console.warn("Failed to parse payload, using current payload");
    }
  }

  // Extract title and body safely
  const title =
    payload?.aps?.alert?.title || payload?.notification?.title || "Hello";
  const body =
    payload?.aps?.alert?.body || payload?.notification?.body || "World";

  return templateLoader.processTemplate(clientCodeTemplate.value, {
    title,
    body,
  });
});

// Generated server code
const serverCode = computed(() => {
  if (isLoadingServerTemplate.value) {
    return "// Loading template...";
  }

  if (!serverCodeTemplate.value) {
    return "// Template not available";
  }

  let payload: any = customPayload.value;

  // If not editing payload here, use the default template or parent payload
  if (!editPayloadHere.value) {
    try {
      // Use APNS payload if available, otherwise use default
      if (props.apnsPayload) {
        payload = JSON.parse(props.apnsPayload);
      } else if (props.fcmPayload) {
        payload = JSON.parse(props.fcmPayload);
      } else {
        payload = JSON.parse(values.PayloadTemplate.DefaultAPNS);
      }
    } catch (e) {
      // Fallback to current payload if parsing fails
      console.warn("Failed to parse payload, using current payload");
    }
  }

  // Extract title and body safely
  const title =
    payload?.aps?.alert?.title || payload?.notification?.title || "Hello";
  const body =
    payload?.aps?.alert?.body || payload?.notification?.body || "World";

  return templateLoader.processTemplate(serverCodeTemplate.value, {
    title,
    body,
  });
});

// Copy to clipboard functions
const copyClientCode = async () => {
  try {
    await navigator.clipboard.writeText(clientCode.value);
    notify("Client code copied to clipboard!", true);
  } catch (err) {
    console.error("Failed to copy client code: ", err);
    notify("Failed to copy client code", false);
  }
};

const copyServerCode = async () => {
  try {
    await navigator.clipboard.writeText(serverCode.value);
    notify("Server code copied to clipboard!", true);
  } catch (err) {
    console.error("Failed to copy server code: ", err);
    notify("Failed to copy server code", false);
  }
};

const resetPayload = () => {
  customPayload.value = defaultPayload;
  customPayloadText.value = JSON.stringify(customPayload.value, null, 2);
};

const updatePayloadFromText = () => {
  try {
    customPayload.value = JSON.parse(customPayloadText.value);
  } catch (e) {
    // Invalid JSON, keep the current payload
    console.warn("Invalid JSON in payload editor");
  }
};

const downloadZip = async () => {
  try {
    // Dynamic import of JSZip
    const JSZip = (await import("jszip")).default;
    const zip = new JSZip();

    // Get current platform and language labels for filenames
    const clientPlatform = clientPlatforms.find(
      (p) => p.value === selectedClientPlatform.value
    );
    const serverLanguage = serverLanguages.find(
      (l) => l.value === selectedServerLanguage.value
    );

    if (!clientPlatform || !serverLanguage) {
      console.error("Invalid platform or language selection");
      return;
    }

    // Determine file extensions based on platform/language
    const getClientFileExtension = (platform: string): string => {
      switch (platform) {
        case "ios-swift":
          return "swift";
        case "ios-objc":
          return "m";
        case "android-kotlin":
          return "kt";
        case "android-java":
          return "java";
        case "flutter":
          return "dart";
        case "react-native":
          return "js";
        default:
          return "txt";
      }
    };

    const getServerFileExtension = (language: string): string => {
      switch (language) {
        case "nodejs":
          return "js";
        case "python":
          return "py";
        case "go":
          return "go";
        case "java":
          return "java";
        case "csharp":
          return "cs";
        case "php":
          return "php";
        default:
          return "txt";
      }
    };

    // Generate filenames
    const clientFileName = `client-${
      selectedClientPlatform.value
    }.${getClientFileExtension(selectedClientPlatform.value)}`;
    const serverFileName = `server-${
      selectedServerLanguage.value
    }.${getServerFileExtension(selectedServerLanguage.value)}`;

    // Add files to zip
    zip.file(clientFileName, clientCode.value);
    zip.file(serverFileName, serverCode.value);

    // Add a README file with information
    const readmeContent = `# Generated Push Notification Code

Generated on: ${new Date().toISOString()}
Platform: ${clientPlatform.label} | Server: ${serverLanguage.label}

## Files Included

### 📱 Client Code
- **${clientFileName}**
- Platform: ${clientPlatform.label}
- Contains client-side code for local notifications

### 🖥️ Server Code  
- **${serverFileName}**
- Language: ${serverLanguage.label}
- Contains server-side code for push notifications via Firebase

${
  selectedServerLanguage.value === "nodejs"
    ? "- **package.json** - Node.js dependencies"
    : ""
}
${
  selectedServerLanguage.value === "python"
    ? "- **requirements.txt** - Python dependencies"
    : ""
}
${
  selectedServerLanguage.value === "go"
    ? "- **go.mod** - Go module definition"
    : ""
}

## 🚀 Quick Start

### Client Setup (${clientPlatform.label})

1. Copy the contents of \`${clientFileName}\` into your ${
      clientPlatform.label
    } project
2. Install required dependencies:
${
  selectedClientPlatform.value === "ios-swift" ||
  selectedClientPlatform.value === "ios-objc"
    ? "   - Import UserNotifications framework\n   - Request notification permissions in your app"
    : selectedClientPlatform.value === "android-kotlin" ||
      selectedClientPlatform.value === "android-java"
    ? "   - Add AndroidX Core library\n   - Add notification permissions to AndroidManifest.xml"
    : selectedClientPlatform.value === "flutter"
    ? "   - Add flutter_local_notifications plugin to pubspec.yaml\n   - Configure platform-specific settings"
    : "   - Install react-native-push-notification\n   - Configure platform-specific settings"
}
3. Configure notification permissions as required by the platform

### Server Setup (${serverLanguage.label})

1. Copy the contents of \`${serverFileName}\` into your ${
      serverLanguage.label
    } project
2. Install dependencies:
${
  selectedServerLanguage.value === "nodejs"
    ? "   ```bash\n   npm install\n   ```"
    : selectedServerLanguage.value === "python"
    ? "   ```bash\n   pip install -r requirements.txt\n   ```"
    : selectedServerLanguage.value === "go"
    ? "   ```bash\n   go mod tidy\n   ```"
    : selectedServerLanguage.value === "java"
    ? "   - Add Firebase Admin SDK to your build.gradle or pom.xml"
    : selectedServerLanguage.value === "csharp"
    ? "   - Install FirebaseAdmin NuGet package"
    : "   ```bash\n   composer require kreait/firebase-php\n   ```"
}
3. Configure Firebase:
   - Download your service account key from Firebase Console
   - Set up Firebase Admin SDK credentials
   - Replace "DEVICE_TOKEN" with actual device tokens

## 🔧 Configuration Notes

### Firebase Setup
1. Go to [Firebase Console](https://console.firebase.google.com)
2. Create a new project or select existing project
3. Navigate to Project Settings > Service Accounts
4. Generate a new private key and download the JSON file
5. Use this file in your server code for authentication

### Testing
- For local notifications: Test on device/simulator directly
- For push notifications: Use Firebase Console or implement the server code
- Always test on physical devices for push notifications

## 📚 Additional Resources
- [Firebase Documentation](https://firebase.google.com/docs)
- [Push Notification Best Practices](https://developer.apple.com/documentation/usernotifications)
- [Android Notification Guide](https://developer.android.com/guide/topics/ui/notifiers/notifications)

---
*Generated by [Remotify](https://github.com/thanhduy26091995/remotify) - Push Notification Testing Tool*
`;

    zip.file("README.md", readmeContent);

    // Add package.json for Node.js projects
    if (selectedServerLanguage.value === "nodejs") {
      const packageJson = {
        name: "push-notification-server",
        version: "1.0.0",
        description: "Generated push notification server code",
        main: serverFileName,
        dependencies: {
          "firebase-admin": "^11.0.0",
        },
        scripts: {
          start: `node ${serverFileName}`,
        },
      };
      zip.file("package.json", JSON.stringify(packageJson, null, 2));
    }

    // Add requirements.txt for Python projects
    if (selectedServerLanguage.value === "python") {
      const requirements = "firebase-admin>=6.0.0\n";
      zip.file("requirements.txt", requirements);
    }

    // Add go.mod for Go projects
    if (selectedServerLanguage.value === "go") {
      const goMod = `module push-notification-server

go 1.19

require (
    firebase.google.com/go v3.13.0+incompatible
)
`;
      zip.file("go.mod", goMod);
    }

    // Generate zip file
    const zipBlob = await zip.generateAsync({ type: "blob" });

    // Create download
    const url = URL.createObjectURL(zipBlob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `remotify-generated-code-${Date.now()}.zip`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);

    // Success feedback (you could add a toast notification here)
    console.log("Code package downloaded successfully");
    notify("Code package downloaded successfully!", true);
  } catch (error) {
    console.error("Failed to generate zip file:", error);
    notify("Failed to generate code package", false);
  }
};
</script>

<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-12">
        <!-- Platform Selection -->
        <div class="row mb-3">
          <div class="col-md-6">
            <label for="clientPlatform" class="form-label"
              >Client Platform:</label
            >
            <select
              id="clientPlatform"
              class="form-select"
              v-model="selectedClientPlatform"
            >
              <option
                v-for="platform in clientPlatforms"
                :key="platform.value"
                :value="platform.value"
              >
                {{ platform.label }}
              </option>
            </select>
          </div>
          <div class="col-md-6">
            <label for="serverLanguage" class="form-label"
              >Server Language:</label
            >
            <select
              id="serverLanguage"
              class="form-select"
              v-model="selectedServerLanguage"
            >
              <option
                v-for="language in serverLanguages"
                :key="language.value"
                :value="language.value"
              >
                {{ language.label }}
              </option>
            </select>
          </div>
        </div>

        <!-- Edit Payload Toggle -->
        <div class="row mb-3">
          <div class="col-12">
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="editPayloadHere"
                v-model="editPayloadHere"
              />
              <label class="form-check-label" for="editPayloadHere">
                Edit Payload Here
              </label>
            </div>
            <small class="text-muted">
              When "Edit Payload Here" is OFF → use payload from APNS/FCM tab<br />
              When ON → show inline JSON editor
            </small>
          </div>
        </div>

        <!-- Payload Editor (conditional) -->
        <div v-if="editPayloadHere" class="row mb-3">
          <div class="col-12">
            <label class="form-label">Payload Editor:</label>
            <div
              class="border rounded p-3"
              :class="isDarkMode ? 'bg-dark text-light' : 'bg-light'"
            >
              <textarea
                class="form-control"
                rows="8"
                v-model="customPayloadText"
                @input="updatePayloadFromText"
                style="font-family: 'Courier New', monospace"
                :class="isDarkMode ? 'bg-dark text-light border-secondary' : ''"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Client Code Section -->
        <div class="row mb-4">
          <div class="col-12">
            <div class="d-flex justify-content-between align-items-center mb-2">
              <h5 class="mb-0">
                Client Code ({{
                  clientPlatforms.find(
                    (p) => p.value === selectedClientPlatform
                  )?.label
                }}):
              </h5>
              <button
                class="btn btn-outline-primary btn-sm"
                @click="copyClientCode"
              >
                <i class="bi bi-clipboard"></i> Copy Client Code
              </button>
            </div>
            <div class="border rounded p-3">
              <pre class="mb-0"><code>{{ clientCode }}</code></pre>
            </div>
          </div>
        </div>

        <!-- Server Code Section -->
        <div class="row mb-4">
          <div class="col-12">
            <div class="d-flex justify-content-between align-items-center mb-2">
              <h5 class="mb-0">
                Server Code ({{
                  serverLanguages.find(
                    (l) => l.value === selectedServerLanguage
                  )?.label
                }}):
              </h5>
              <button
                class="btn btn-outline-primary btn-sm"
                @click="copyServerCode"
              >
                <i class="bi bi-clipboard"></i> Copy Server Code
              </button>
            </div>
            <div class="border rounded p-3">
              <pre class="mb-0"><code>{{ serverCode }}</code></pre>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="mb-4 d-flex justify-content-between">
          <div class="d-flex justify-content-between">
            <button @click="resetPayload" type="button" class="btn btn-danger">
              Reset <i class="bi bi-eraser text-white"></i>
            </button>
            &nbsp;&nbsp;
            <button @click="downloadZip" type="button" class="btn btn-primary">
              Download <i class="bi bi-arrow-down-circle"></i>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-size: 0.9rem;
  max-height: 400px;
  overflow-y: auto;
}

code {
  font-family: "Courier New", Monaco, monospace;
}

.form-select,
.form-control {
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.btn {
  transition: all 0.15s ease-in-out;
}
</style>
