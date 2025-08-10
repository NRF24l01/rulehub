<template>
    <nav class="bg-white shadow px-3 py-4">
        <div class="container mx-auto flex items-center justify-between">
            <!-- Logo and Brand -->
            <router-link to="/" class="flex items-center space-x-4">
                <img src="/ico.png" alt="Logo" class="h-10 w-10" />
                <span class="text-2xl font-bold text-gray-800">Rule Hub</span>
            </router-link>

            <!-- Burger button (mobile only) -->
            <button
                class="md:hidden flex flex-col justify-center items-center w-10 h-10 focus:outline-none"
                @click="isMenuOpen = !isMenuOpen"
                aria-label="Toggle navigation"
            >
                <span :class="['block h-0.5 w-6 bg-gray-800 transition-all duration-300', isMenuOpen ? 'rotate-45 translate-y-2' : '']"></span>
                <span :class="['block h-0.5 w-6 bg-gray-800 my-1 transition-all duration-300', isMenuOpen ? 'opacity-0' : '']"></span>
                <span :class="['block h-0.5 w-6 bg-gray-800 transition-all duration-300', isMenuOpen ? '-rotate-45 -translate-y-2' : '']"></span>
            </button>

            <!-- Navigation Links -->
            <div
                class="flex items-center space-x-8 md:flex"
                :class="[
                    'md:static md:flex-row md:space-x-8',
                    isMenuOpen ? 'absolute top-16 left-0 w-full bg-white flex flex-col space-y-4 px-6 py-4 shadow-lg z-50' : 'hidden',
                    'md:!flex'
                ]"
            >
                <router-link 
                    v-if="localIsAuthenticated" 
                    to="/create" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                    @click="closeMenuOnMobile"
                >
                    <i class="fa fa-plus-circle mr-2" aria-hidden="true"></i>
                    Create Rule
                </router-link>
                
                <router-link 
                    to="/about" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                    @click="closeMenuOnMobile"
                >
                    <i class="fa fa-info-circle mr-2" aria-hidden="true"></i>
                    About
                </router-link>
                
                <a 
                    href="https://github.com/NRF24l01/rulehub" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    class="flex items-center text-lg text-gray-600 hover:text-blue-600 font-medium transition duration-200"
                    @click="closeMenuOnMobile"
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

const isMenuOpen = ref(false);

function closeMenuOnMobile() {
  isMenuOpen.value = false;
}

onMounted(async () => {
  if (!authStore.isAuthenticated && authStore.accessToken) {
    const refreshed = await authStore.refreshToken();
    localIsAuthenticated.value = authStore.isAuthenticated;
  } else {
    localIsAuthenticated.value = authStore.isAuthenticated;
  }
});
</script>

<style scoped>
@media (min-width: 768px) {
  .md\\:hidden { display: none !important; }
  .md\\:flex { display: flex !important; }
  .md\\:static { position: static !important; }
  .md\\:flex-row { flex-direction: row !important; }
  .md\\:space-x-8 > :not([hidden]) ~ :not([hidden]) { margin-left: 2rem !important; margin-top: 0 !important; }
}
</style>