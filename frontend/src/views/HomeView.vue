<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white overflow-hidden relative">
      <div class="absolute opacity-10 right-0 top-0 w-1/2 h-full">
        <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg" class="w-full h-full">
          <path fill="white" d="M42.7,-62.9C54.9,-53.5,64,-39.7,70.1,-24.5C76.3,-9.3,79.5,7.3,75.5,22.1C71.5,37,60.3,50.1,46.6,58.7C32.9,67.3,16.5,71.4,0.2,71.1C-16,70.8,-32,66.2,-45.8,57.3C-59.6,48.3,-71.2,35.1,-76.8,19.3C-82.4,3.5,-82,-14.8,-74.9,-29.8C-67.8,-44.7,-54,-56.2,-39.6,-64.4C-25.2,-72.6,-10.1,-77.3,3.4,-82C16.9,-86.6,33.9,-91.1,41.8,-84.1C49.7,-77,60.5,-58.3,42.7,-62.9Z" transform="translate(100 100)" />
        </svg>
      </div>
      <div class="container mx-auto px-4 py-20 md:py-28 relative z-10">
        <div class="flex flex-col md:flex-row items-center">
          <div class="md:w-1/2 mb-8 md:mb-0" data-aos="fade-right">
            <h1 class="text-4xl md:text-6xl font-bold mb-6 leading-tight">Управляйте правилами вашего проекта</h1>
            <p class="text-xl mb-8 text-blue-100">Простой и удобный способ создания, редактирования и публикации правил для вашего сообщества или организации.</p>
            <div class="flex flex-col sm:flex-row gap-4">
              <router-link to="/editor" class="bg-white text-indigo-700 hover:bg-indigo-50 px-8 py-4 rounded-lg font-medium text-center transition duration-300 transform hover:scale-105 shadow-lg">
                Создать правила
              </router-link>
              <router-link to="/rules" class="bg-transparent border-2 border-white hover:bg-white hover:text-indigo-700 px-8 py-4 rounded-lg font-medium text-center transition duration-300">
                Просмотреть примеры
              </router-link>
            </div>
          </div>
          <div class="md:w-1/2" data-aos="fade-left" data-aos-delay="200">
            <img src="/doc.svg" alt="Rules Editor" class="w-full max-w-md mx-auto transform hover:rotate-2 transition-transform duration-300 drop-shadow-2xl" />
          </div>
        </div>
      </div>
    </section>

    <!-- Stats Counter -->
    <section class="bg-white py-8 shadow-md relative z-10">
      <div class="container mx-auto px-4">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
          <div class="p-4">
            <div class="text-3xl font-bold text-indigo-600 mb-1">{{ stats.projects }}+</div>
            <div class="text-gray-600">Страничек на сайте</div>
          </div>
          <div class="p-4">
            <div class="text-3xl font-bold text-indigo-600 mb-1">{{ stats.documents }}k+</div>
            <div class="text-gray-600">Завайбленых документов</div>
          </div>
          <div class="p-4">
            <div class="text-3xl font-bold text-indigo-600 mb-1">{{ stats.clients }}%</div>
            <div class="text-gray-600">Довольных уток</div>
          </div>
          <div class="p-4">
            <div class="text-3xl font-bold text-indigo-600 mb-1">{{ stats.support }}</div>
            <div class="text-gray-600">Техподдержка</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-4" data-aos="fade-up">Возможности платформы</h2>
        <p class="text-gray-600 text-center mb-12 max-w-2xl mx-auto" data-aos="fade-up" data-aos-delay="100">Наш сервис предлагает все необходимые инструменты для эффективного управления правилами вашего проекта.</p>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="bg-white p-8 rounded-xl shadow-md hover:shadow-xl transition duration-300 transform hover:-translate-y-2" data-aos="fade-up" data-aos-delay="150">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold mb-3">Встроенный редактор</h3>
            <p class="text-gray-600">Удобный редактор с поддержкой форматирования для создания четких и понятных правил.</p>
          </div>
          
          <div class="bg-white p-8 rounded-xl shadow-md hover:shadow-xl transition duration-300 transform hover:-translate-y-2" data-aos="fade-up" data-aos-delay="300">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold mb-3">Адаптивный дизайн</h3>
            <p class="text-gray-600">Работайте с правилами на любом устройстве - от смартфона до настольного компьютера.</p>
          </div>
          
          <div class="bg-white p-8 rounded-xl shadow-md hover:shadow-xl transition duration-300 transform hover:-translate-y-2" data-aos="fade-up" data-aos-delay="450">
            <div class="text-indigo-600 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold mb-3">Простое развертывание</h3>
            <p class="text-gray-600">Запустите свой сервис правил в один клик с помощью Docker-контейнера.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="bg-gray-100 py-20">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-4" data-aos="fade-up">Как это работает</h2>
        <p class="text-gray-600 text-center mb-12 max-w-2xl mx-auto" data-aos="fade-up" data-aos-delay="100">Четыре простых шага для запуска вашей системы правил</p>
        
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div class="flex flex-col items-center text-center group" data-aos="fade-up" data-aos-delay="150">
            <div class="bg-indigo-600 text-white rounded-full w-16 h-16 flex items-center justify-center text-2xl font-bold mb-5 transform group-hover:scale-110 transition-transform duration-300 shadow-lg">1</div>
            <h3 class="text-xl font-semibold mb-3">Разверните проект</h3>
            <p class="text-gray-600">Используйте Docker для быстрого запуска сервиса правил.</p>
          </div>
          
          <div class="flex flex-col items-center text-center group" data-aos="fade-up" data-aos-delay="300">
            <div class="bg-indigo-600 text-white rounded-full w-16 h-16 flex items-center justify-center text-2xl font-bold mb-5 transform group-hover:scale-110 transition-transform duration-300 shadow-lg">2</div>
            <h3 class="text-xl font-semibold mb-3">Создайте правила</h3>
            <p class="text-gray-600">Используйте встроенный редактор для составления правил.</p>
          </div>
          
          <div class="flex flex-col items-center text-center group" data-aos="fade-up" data-aos-delay="450">
            <div class="bg-indigo-600 text-white rounded-full w-16 h-16 flex items-center justify-center text-2xl font-bold mb-5 transform group-hover:scale-110 transition-transform duration-300 shadow-lg">3</div>
            <h3 class="text-xl font-semibold mb-3">Опубликуйте</h3>
            <p class="text-gray-600">Сохраните и опубликуйте правила для вашей аудитории.</p>
          </div>
          
          <div class="flex flex-col items-center text-center group" data-aos="fade-up" data-aos-delay="600">
            <div class="bg-indigo-600 text-white rounded-full w-16 h-16 flex items-center justify-center text-2xl font-bold mb-5 transform group-hover:scale-110 transition-transform duration-300 shadow-lg">4</div>
            <h3 class="text-xl font-semibold mb-3">Обновляйте</h3>
            <p class="text-gray-600">Легко вносите изменения и поддерживайте актуальность.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Testimonials -->
    <section class="py-20">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-4" data-aos="fade-up">Отзывы клиентов</h2>
        <p class="text-gray-600 text-center mb-12 max-w-2xl mx-auto" data-aos="fade-up" data-aos-delay="100">Что говорят наши пользователи о платформе</p>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="bg-white p-6 rounded-xl shadow-md" data-aos="fade-up" data-aos-delay="150">
            <div class="flex items-center mb-4">
              <div class="text-yellow-500 flex">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" v-for="i in 5" :key="i">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </div>
            </div>
            <p class="text-gray-600 mb-4">"Этот сервис полностью изменил способ, которым мы управляем правилами нашего сообщества. Все стало намного проще и эффективнее."</p>
            <div class="font-semibold">Стас Админов</div>
            <div class="text-gray-500 text-sm">Модератор.</div>
          </div>
          
          <div class="bg-white p-6 rounded-xl shadow-md" data-aos="fade-up" data-aos-delay="300">
            <div class="flex items-center mb-4">
              <div class="text-yellow-500 flex">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" v-for="i in 5" :key="i">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </div>
            </div>
            <p class="text-gray-600 mb-4">"Удобный интерфейс и отличная адаптивность для мобильных устройств. Мы развернули сервис за считанные минуты."</p>
            <div class="font-semibold">Стас Девопсов</div>
            <div class="text-gray-500 text-sm">Единственный кто может общаться с кластером на балконе</div>
          </div>
          
          <div class="bg-white p-6 rounded-xl shadow-md" data-aos="fade-up" data-aos-delay="450">
            <div class="flex items-center mb-4">
              <div class="text-yellow-500 flex">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" v-for="i in 5" :key="i">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </div>
            </div>
            <p class="text-gray-600 mb-4">"Превосходное решение! Больше не нужно беспокоиться о форматировании и доступе к правилам - всё работает из коробки. Игроков уже можно ограничивать!"</p>
            <div class="font-semibold">Стас Козлов</div>
            <div class="text-gray-500 text-sm">Владелец игрового сервера</div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 bg-gradient-to-r from-indigo-600 to-blue-500 text-white">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl md:text-4xl font-bold mb-6" data-aos="fade-up">Готовы начать?</h2>
        <p class="text-xl text-blue-100 mb-8 max-w-2xl mx-auto" data-aos="fade-up" data-aos-delay="100">Создайте свои первые правила прямо сейчас или изучите документацию, чтобы узнать больше о возможностях платформы.</p>
        <div class="flex flex-col sm:flex-row justify-center gap-4" data-aos="fade-up" data-aos-delay="200">
          <router-link to="/editor" class="bg-white hover:bg-gray-100 text-indigo-600 px-8 py-4 rounded-lg font-medium transition duration-300 shadow-lg transform hover:scale-105">
            Начать работу
          </router-link>
          <router-link to="/docs" class="bg-transparent border-2 border-white hover:bg-white hover:text-indigo-600 px-8 py-4 rounded-lg font-medium transition duration-300">
            Документация
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue';
import AOS from 'aos';
import 'aos/dist/aos.css';

// Статистика для анимации
const stats = ref({
  projects: 0,
  documents: 0,
  clients: 0,
  support: '24/7'
});

// Целевые значения
const targets = {
  projects: 12,
  documents: 12,
  clients: 100, // clients будет расти бесконечно, но начальное значение для других целей
  support: '0/0'
};

// Анимация чисел и support
function animateStats() {
  let p = 0;
  let d = 0;
  let c = 0;
  let sHours = 24;
  let sDays = 7;
  const pTarget = targets.projects;
  const dTarget = targets.documents;
  const duration = 2000; // ms
  const steps = 50;
  const pStep = Math.ceil(pTarget / steps);
  const dStep = Math.ceil(dTarget / steps);
  const sHoursStep = Math.ceil(sHours / steps);
  const sDaysStep = Math.ceil(sDays / steps);

  // clients будет расти бесконечно
  const cStep = 2; // скорость роста клиентов

  const interval = setInterval(() => {
    if (p < pTarget) p = Math.min(p + pStep, pTarget);
    if (d < dTarget) d = Math.min(d + dStep, dTarget);

    // clients растет бесконечно
    c += cStep;

    // support падение
    if (sHours > 0) sHours = Math.max(sHours - sHoursStep, 0);
    if (sDays > 0) sDays = Math.max(sDays - sDaysStep, 0);

    stats.value.projects = p;
    stats.value.documents = d;
    stats.value.clients = c;
    stats.value.support = `${sHours}/${sDays}`;

    // Когда все закончено (кроме clients)
    if (p === pTarget && d === dTarget && sHours === 0 && sDays === 0) {
      stats.value.support = targets.support;
      // Не останавливаем интервал, чтобы clients продолжал расти
    }
  }, duration / steps);
}

onMounted(() => {
  AOS.init({
    duration: 1000,
    easing: 'ease-in-out',
    once: true
  });
  animateStats();
});
</script>
