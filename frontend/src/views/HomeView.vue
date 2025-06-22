<script setup>
import { ref, onMounted } from 'vue'

// --- State Management ---
// Reactive variable for the list of files from the API
const files = ref([])
// To hold any error messages during fetch
const error = ref(null)
// To hold the file selected by the user in the form
const selectedFile = ref(null)
// To show feedback during and after the upload
const uploadStatus = ref('')

// --- API Logic ---
// We've extracted the fetching logic into its own reusable function
async function fetchFiles() {
  try {
    error.value = null // Reset error on new fetch
    const response = await fetch('http://localhost:8080/api/v1/files')
    if (!response.ok) {
      throw new Error('Network response was not ok')
    }
    const data = await response.json()
    files.value = data
  } catch (e) {
    error.value = e.message
    console.error('There was an error fetching the files:', e)
  }
}

// Function to handle the form submission
async function handleUpload() {
  // 1. Check if a file is actually selected
  if (!selectedFile.value) {
    uploadStatus.value = 'Please select a file first.'
    return
  }

  // 2. Use the FormData API to prepare the file for sending.
  // This is the standard way to send files via HTTP.
  const formData = new FormData()
  // The key 'file' must match what our Go backend expects: r.FormFile("file")
  formData.append('file', selectedFile.value)

  uploadStatus.value = 'Uploading...'

  try {
    // 3. Make the POST request using fetch
    const response = await fetch('http://localhost:8080/api/v1/files', {
      method: 'POST',
      body: formData,
      // NOTE: We DO NOT set the 'Content-Type' header ourselves.
      // The browser will automatically set it to 'multipart/form-data'
      // with the correct boundary when we use FormData.
    })

    if (!response.ok) {
      throw new Error('Upload failed.')
    }

    uploadStatus.value = 'Upload successful!'

    // 4. Refresh the file list to show the new file!
    await fetchFiles()

  } catch (e) {
    uploadStatus.value = `Error: ${e.message}`
    console.error('Error uploading file:', e)
  }
}

// Function to update our 'selectedFile' ref when the user chooses a file
function handleFileChange(event) {
  selectedFile.value = event.target.files[0]
}

// --- Lifecycle Hook ---
// Fetch the initial list of files when the component is first mounted
onMounted(fetchFiles)

</script>

<template>
  <div>
    <h1>StorageCage</h1>

    <section class="upload-section">
      <h2>Upload New File</h2>
      <form @submit.prevent="handleUpload">
        <input type="file" @change="handleFileChange" required />
        <button type="submit">Upload</button>
      </form>
      <p v-if="uploadStatus" class="status-message">{{ uploadStatus }}</p>
    </section>

    <section class="file-list-section">
      <h2>Stored Files</h2>
      <div v-if="error" class="error">
        <p>Error fetching files: {{ error }}</p>
      </div>
      <ul v-else-if="files && files.length">
        <li v-for="file in files" :key="file.id">
          <span>{{ file.name }}</span>
          <span class="size">({{ file.size }} bytes)</span>
        </li>
      </ul>
      <div v-else>
        <p>No files found or still loading...</p>
      </div>
    </section>

  </div>
</template>

<style scoped>
h1 {
  border-bottom: 2px solid #eee;
  padding-bottom: 0.5rem;
}

.upload-section, .file-list-section {
  margin-top: 2rem;
}

.status-message {
  margin-top: 1rem;
  font-style: italic;
  color: #333;
}

.error {
  color: red;
}

.size {
  margin-left: 1rem;
  color: #666;
  font-size: 0.9em;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>