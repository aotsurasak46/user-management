<script setup>
import { useRouter } from 'vue-router'
import SidebarLink from './SidebarLink.vue';
import { useUserAccountStore } from '../stores/userAccount'
import { ref } from 'vue'

const userAccountStore = useUserAccountStore()
const router = useRouter()
const isMobileMenuOpen = ref(false)

const handleLogout = async () => {
  await userAccountStore.logout()
  router.push({name : 'login'})
}

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}
</script>

<template>
  <!-- Mobile sidebar toggle button-->
  <button
    @click="toggleMobileMenu"
    class="md:hidden fixed top-4 left-4 bg-gray-800 text-white p-3 rounded-md z-51 text-xs"
  >
    <i class="pi pi-bars"></i>
  </button>

  <!-- Mobile sidebar -->
  <div v-if="isMobileMenuOpen" class="fixed top-0 left-0 right-0 bottom-0 z-50 bg-black/20 flex items-center ">
    <aside class="md:hidden  h-screen border-r border-gray-300 w-17 bg-white">
      <div class="flex flex-col justify-between h-full">
        <div class="flex flex-col items-center p-3">
          <div class='h-15'></div>
          <div class="w-full flex flex-col justify-start">
            <span class="text-[10px] text-[#abadab] mb-1 hidden lg:inline ">USER MANAGEMENT</span>
            <SidebarLink route-name="user-list" :icon="'pi pi-user'" label="User List"/>
            <SidebarLink route-name="chat" :icon="'pi pi-comments'" label="Chat"/>
          </div>
        </div>
        <div class="flex flex-col items-center justify-center lg:border-t border-t-gray-300 p-4 bg-light ">
          <div class="flex w-10 h-10 rounded-full bg-gray-200  items-center justify-center text-gray-500 mb-4 lg:mb-0">
            <span>{{ userAccountStore.user?.name ? userAccountStore.user.name.charAt(0).toUpperCase() : '' }}</span>
          </div>
          <button @click="handleLogout" class="cursor-pointer hover:bg-gray-200 rounded-sm disabled:cursor-auto disabled:hover:bg-black p-2 shadow-md">
              <i class="pi pi-sign-out"></i>
          </button>
        </div> 
      </div>
    </aside>
  </div>

  <!-- Window sidebar -->
  <aside class="hidden md:flex flex-col justify-between h-screen border-r border-gray-300 w-20 lg:w-50  ">
    <div class="flex flex-col items-center p-3">
      <div class="w-full flex flex-col justify-start">
        <span class="text-[10px] text-[#abadab] mb-1 hidden lg:inline ">USER MANAGEMENT</span>
        <SidebarLink route-name="user-list" :icon="'pi pi-user'" label="User List"/>
        <SidebarLink route-name="chat" :icon="'pi pi-comments'" label="Chat"/>
      </div>
    </div>
    <div class="flex flex-col lg:flex-row  items-center justify-center lg:border-t border-t-gray-300 p-4 bg-light ">
      <div class="hidden md:flex w-10 h-10 rounded-full bg-gray-200  items-center justify-center text-gray-500 mb-4 lg:mb-0">
        <span>{{ userAccountStore.user?.name ? userAccountStore.user.name.charAt(0).toUpperCase() : '' }}</span>
      </div>
      <div class="hidden lg:flex flex-col justify-center px-4 min-w-0 flex-1 overflow-hidden">
        <span class="text-sm text-gray-900 font-medium truncate">{{ userAccountStore.user?.name ? userAccountStore.user.name : '' }}</span>
        <span class="text-xs text-gray-500 font-normal truncate">{{ userAccountStore.user?.email ? userAccountStore.user.email : '' }}</span>
      </div>
      <button @click="handleLogout" class="cursor-pointer hover:bg-gray-200 rounded-sm disabled:cursor-auto disabled:hover:bg-black p-2 shadow-md">
          <i class="pi pi-sign-out"></i>
      </button>
    </div> 
  </aside>

</template>