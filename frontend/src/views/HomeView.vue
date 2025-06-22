<script setup>
import { ref, onMounted } from 'vue'

// Create a "ref", which is a reactive variable to hold our data.
// When this variable changes, Vue will automatically update the template.
const files = ref([])
const error = ref(null)

// 'onMounted' is a lifecycle hook. The code inside it will run
// automatically when the component is first added to the page.
onMounted(async () => {
  try {
    // Use the browser's built-in 'fetch' API to make a GET request to our backend.
    const response = await fetch('http://localhost:8080/api/v1/files')

    if (!response.ok) {
      throw new Error('Network response was not ok')
    }

    // Parse the JSON response from the backend.
    const data = await response.json()

    // Update our reactive variable with the data from the API.
    files.value = data
  } catch (e) {
    // If an error occurs, store it so we can display it to the user.
    error.value = e.message
    console.error('There was an error fetching the files:', e)
  }
})
</script>

<template>
  <div>
    <h1>StorageCage Files</h1>

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
  </div>
</template>

<style scoped>
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
}
</style>