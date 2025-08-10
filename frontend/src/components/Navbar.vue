<template>
    <nav class="bg-white shadow px-3 py-4">
        <div class="container mx-auto flex items-center justify-between">
            <!-- Logo and Brand -->
            <router-link to="/" class="flex items-center space-x-4">
                <img src="/ico.png" alt="Logo" class="h-10 w-10" />
                <span class="text-2xl font-bold text-gray-800">Rule Hub</span>
            </router-link>

            <!-- Navigation Links -->
            <div class="flex items-center space-x-8">
                <router-link 
                    v-if="localIsAuthenticated" 
                    to="/create" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                >
                    <i class="fa fa-plus-circle mr-2" aria-hidden="true"></i>
                    Create Rule
                </router-link>
                
                <router-link 
                    to="/about" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                >
                    <i class="fa fa-info-circle mr-2" aria-hidden="true"></i>
                    About
                </router-link>
                
                <a 
                    href="https://github.com/NRF24l01/rulehub" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                >
                    <i class="fab fa-github mr-2" aria-hidden="true"></i>
                    GitHub
                </a>
            </div>
        </div>
    </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const localIsAuthenticated = ref(authStore.isAuthenticated);

onMounted(async () => {
  if (!authStore.isAuthenticated && authStore.accessToken) {
    const refreshed = await authStore.refreshToken();
    localIsAuthenticated.value = authStore.isAuthenticated;
  } else {
    localIsAuthenticated.value = authStore.isAuthenticated;
  }
});
</script>