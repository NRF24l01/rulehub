<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
      <div class="container mx-auto px-4 py-16 md:py-24">
        <div class="flex flex-col md:flex-row items-center">
          <div class="md:w-1/2 mb-8 md:mb-0">
            <h1 class="text-4xl md:text-5xl font-bold mb-4">Управляйте правилами вашего проекта</h1>
            <p class="text-xl mb-8">Простой и удобный способ создания, редактирования и публикации правил для вашего сообщества или организации.</p>
            <div class="flex flex-col sm:flex-row gap-4">
              <router-link to="/editor" class="bg-white text-indigo-700 hover:bg-indigo-50 px-6 py-3 rounded-lg font-medium text-center transition duration-300">
                Создать правила
              </router-link>
              <router-link to="/rules" class="bg-transparent border border-white hover:bg-white hover:text-indigo-700 px-6 py-3 rounded-lg font-medium text-center transition duration-300">
                Просмотреть примеры
              </router-link>
            </div>
          </div>
          <div class="md:w-1/2">
            <img src="/doc.svg" alt="Rules Editor" class="w-full max-w-md mx-auto" />
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-12">Возможности платформы</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition duration-300">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Встроенный редактор</h3>
            <p class="text-gray-600">Удобный редактор с поддержкой форматирования для создания четких и понятных правил.</p>
          </div>
          
          <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition duration-300">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Адаптивный дизайн</h3>
            <p class="text-gray-600">Работайте с правилами на любом устройстве - от смартфона до настольного компьютера.</p>
          </div>
          
          <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition duration-300">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Простое развертывание</h3>
            <p class="text-gray-600">Запустите свой сервис правил в один клик с помощью Docker-контейнера.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="bg-gray-100 py-16">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-12">Как это работает</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div class="flex flex-col items-center text-center">
            <div class="bg-indigo-600 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">1</div>
            <h3 class="text-xl font-semibold mb-2">Разверните проект</h3>
            <p class="text-gray-600">Используйте Docker для быстрого запуска сервиса правил.</p>
          </div>
          
          <div class="flex flex-col items-center text-center">
            <div class="bg-indigo-600 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">2</div>
            <h3 class="text-xl font-semibold mb-2">Создайте правила</h3>
            <p class="text-gray-600">Используйте встроенный редактор для составления правил.</p>
          </div>
          
          <div class="flex flex-col items-center text-center">
            <div class="bg-indigo-600 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">3</div>
            <h3 class="text-xl font-semibold mb-2">Опубликуйте</h3>
            <p class="text-gray-600">Сохраните и опубликуйте правила для вашей аудитории.</p>
          </div>
          
          <div class="flex flex-col items-center text-center">
            <div class="bg-indigo-600 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">4</div>
            <h3 class="text-xl font-semibold mb-2">Обновляйте</h3>
            <p class="text-gray-600">Легко вносите изменения и поддерживайте актуальность.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-16">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl md:text-4xl font-bold mb-6">Готовы начать?</h2>
        <p class="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">Создайте свои первые правила прямо сейчас или изучите документацию, чтобы узнать больше о возможностях платформы.</p>
        <div class="flex flex-col sm:flex-row justify-center gap-4">
          <router-link to="/editor" class="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-3 rounded-lg font-medium transition duration-300">
            Начать работу
          </router-link>
          <router-link to="/docs" class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-8 py-3 rounded-lg font-medium transition duration-300">
            Документация
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'HomeView',
  components: {},
  setup() {
    return {};
  }
}
</script>
