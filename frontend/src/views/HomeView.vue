<script setup>
import { ref, onMounted } from 'vue'

// --- State Management ---
const files = ref([])
const error = ref(null)
const selectedFile = ref(null)
const uploadStatus = ref('')

// --- API Logic ---
async function fetchFiles() {
  try {
    error.value = null
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

async function handleUpload() {
  if (!selectedFile.value) {
    uploadStatus.value = 'Please select a file first.'
    return
  }
  const formData = new FormData()
  formData.append('file', selectedFile.value)
  uploadStatus.value = 'Uploading...'
  try {
    const response = await fetch('http://localhost:8080/api/v1/files', {
      method: 'POST',
      body: formData,
    })
    if (!response.ok) {
      throw new Error('Upload failed.')
    }
    uploadStatus.value = 'Upload successful!'
    await fetchFiles() // Refresh the list
  } catch (e) {
    uploadStatus.value = `Error: ${e.message}`
    console.error('Error uploading file:', e)
  }
}

// --- NEW: Delete Functionality ---
async function handleDelete(fileName) {
  // 1. Ask for user confirmation. This is crucial for destructive actions!
  if (!confirm(`Are you sure you want to delete "${fileName}"?`)) {
    return
  }

  try {
    error.value = null
    // 2. Make the DELETE request with the filename in the URL
    const response = await fetch(`http://localhost:8080/api/v1/files/${fileName}`, {
      method: 'DELETE',
    })

    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || 'Could not delete file.')
    }
    
    // 3. If successful, refresh the file list to reflect the change
    await fetchFiles()

  } catch (e) {
    error.value = e.message
    console.error('Error deleting file:', e)
  }
}


function handleFileChange(event) {
  selectedFile.value = event.target.files[0]
}

// --- Lifecycle Hook ---
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
        <p>Error: {{ error }}</p>
      </div>
      <ul v-else-if="files && files.length">
        <li v-for="file in files" :key="file.id">
          <div>
            <span>{{ file.name }}</span>
            <span class="size">({{ file.size }} bytes)</span>
          </div>
          <button @click="handleDelete(file.name)" class="delete-btn">Delete</button>
        </li>
      </ul>
      <div v-else>
        <p>No files found or still loading...</p>
      </div>
    </section>

  </div>
</template>

<style scoped>
/* ... (existing styles) ... */
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
  color: #d32f2f;
  background-color: #ffcdd2;
  padding: 1rem;
  border-radius: 4px;
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
  padding: 0.5rem 0.25rem;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.delete-btn {
  background-color: #d32f2f;
  color: white;
  border: none;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  cursor: pointer;
}
.delete-btn:hover {
  background-color: #c62828;
}
</style>