<script setup>
import { onMounted, ref, computed} from 'vue';
import { useUserStore } from '../stores/user'
import { useUserAccountStore } from '../stores/userAccount'
import UserUpdateModal from '@/components/UserUpdateModal.vue';
import ConfirmDeleteModal from '../components/ConfirmDeleteModal.vue'
import { useChatStore } from '@/stores/chat';
import { useRouter } from 'vue-router';
const userStore = useUserStore()
const userAccoutStore = useUserAccountStore()
const chatStore = useChatStore()
const router = useRouter()

const isLoading = ref(false)

const isShowUpdateModal = ref(false)
const isShowConfirmDeleteModal = ref(false)
const selectedUserId = ref(null)
const isEditMode = ref(false)
const searchValue = ref('')
const roleFilter = ref('all')
const sortBy = ref('')
const filteredUsers = computed(() => {
  let filtered = [...userStore.users]

  if (roleFilter.value !== 'all') {
    filtered = filtered.filter(user => user.role === roleFilter.value)
  }

  if (searchValue.value.trim()) {
    const search = searchValue.value.toLowerCase()
    filtered = filtered.filter(user =>
      user.name.toLowerCase().includes(search) || user.email.toLowerCase().includes(search)
    )
  }
  
  if (sortBy.value === 'name') {
    filtered.sort((a, b) => a.name.localeCompare(b.name))
  } else if (sortBy.value === 'email') {
      filtered.sort((a, b) => a.email.localeCompare(b.email))
  } else if (sortBy.value === 'created-newest') {
    filtered.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt))
  } else if (sortBy.value === 'created-oldest') {
    filtered.sort((a, b) => new Date(a.CreatedAt) - new Date(b.CreatedAt))
  } else if (sortBy.value === 'updated-newest') {
    filtered.sort((a, b) => new Date(b.UpdatedAt) - new Date(a.UpdatedAt))
  } else if (sortBy.value === 'updated-oldest') {
    filtered.sort((a, b) => new Date(a.UpdatedAt) - new Date(b.UpdatedAt))
  }
  return filtered
})

const fetchUsers = async () => {
  try {
    isLoading.value = true
    await userStore.loadUsers()
  } catch (error) {
    toast.error(error.message, {
        autoClose: 3000,
        position: 'top-right',
    })
    console.log('Create user failed', error)
  }finally{
    isLoading.value = false
  }
  
}

onMounted(() => {
  fetchUsers()
})

const toggleUpdateModal = (userId = null, edit = false, shouldRefresh = false) => {
  isShowUpdateModal.value = !isShowUpdateModal.value
  selectedUserId.value = userId
  isEditMode.value = edit
  if (shouldRefresh) {
    fetchUsers()
  }
}

const toggleConfirmDeleteModal = (userId = null ,shouldRefresh = false) => {
  isShowConfirmDeleteModal.value = !isShowConfirmDeleteModal.value
  selectedUserId.value = userId
  if (shouldRefresh) {
    fetchUsers()
  }
}

const sendMessage = (user) => {
  chatStore.setSelectedChatUser(user)
  router.push({name : "chat"})
}

</script>

<template>
  <div class="pl-5 pt-5">
  <div class="flex flex-col">
    <div class="w-full">
      <div class="flex justify-end items-center mb-4 pr-5 md:justify-between">
        <h2 class="text-xl font-semibold hidden md:inline">All users <span class="text-gray-500 ">{{ filteredUsers.length }}</span></h2>
        <div class="flex gap-2 overflow-x-auto whitespace-nowrap">
            <div class="flex items-center border px-3 py-1 rounded-md w-25 md:w-full">
              <input type="text" placeholder="Search" v-model="searchValue" class="w-full outline-none truncate overflow-hidden text-ellipsis whitespace-nowrap">
              <i class="pi pi-search ml-2 text-gray-500"></i>
            </div>
          <select v-model="roleFilter" class="border px-2 py-1 rounded hover:bg-gray-100">
            <option value="all">All Roles</option>
            <option value="admin">Admin</option>
            <option value="user">User</option>
          </select>
          <select v-model="sortBy" class="border px-2 py-1 rounded w-35 md:w-full truncate overflow-hidden text-ellipsis whitespace-nowrap">
            <option value="">Sort By</option>
            <option value="name">Name</option>
            <option value="email">Email</option>
            <option value="created-newest">Created: Newest</option>
            <option value="created-oldest">Created: Oldest</option>
            <option value="updated-newest">Updated: Newest</option>
            <option value="updated-oldest">Updated: Oldest</option>
          </select>
          <button v-if="userAccoutStore.user.role === 'admin' " @click="toggleUpdateModal(selectedUserId, false)" :disabled="userAccoutStore.user.role !== 'admin'" class="flex justify-center items-center sm:block w-8 sm:w-full bg-black hover:bg-gray-600 text-white px-4 py-1 rounded-md cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed">
            <i class="pi pi-plus"></i>
            <span class="hidden lg:inline px-2">Add User</span>
          </button>
        </div>
      </div>
    </div>
  </div>
  <div class="max-h-[90vh] overflow-y-auto
      [&::-webkit-scrollbar]:w-2
      [&::-webkit-scrollbar-track]:rounded-full
    [&::-webkit-scrollbar-track]:bg-gray-100
      [&::-webkit-scrollbar-thumb]:rounded-full
    [&::-webkit-scrollbar-thumb]:bg-gray-300">
  <table class="w-full text-left table-auto ">
        <thead class="text-left sticky top-0 bg-white text-sm text-gray-500 border-b border-b-gray-400">
          <tr>
            <th class="py-2 ">User name</th>
            <th class="py-2 ">Access</th>
            <th class="py-2 hidden lg:table-cell">Last updated</th>
            <th class="py-2 hidden lg:table-cell">Date created</th>
            <th class="py-2 ">Action</th>
          </tr>
        </thead>
        <tbody v-if="isLoading">
          <tr>
            <td colspan="5" class="text-center py-4">
              <div class="flex justify-center items-center space-x-2">
                <span>Loading...</span>
              </div>
            </td>
          </tr>
        </tbody>
        <tbody v-else v-for="user in filteredUsers" :key="user.ID" class="text-sm text-gray-700">
          <tr class="border-b-gray-400 hover:bg-gray-50">
            <td class="py-3 flex ">
              <div class="min-w-10 h-10 rounded-full mr-3 bg-gray-200 flex items-center justify-center text-gray-500">
                <span>{{ user.name.charAt(0).toUpperCase() }}</span>
              </div>
              <div>
                <div class="font-medium text-gray-900"><span>{{ user.name }}</span></div>
                <div class="text-xs text-gray-500"><span>{{ user.email }}</span></div>
              </div>
              <div v-if="user.ID === userAccoutStore.user.id" class="items-center px-3 hidden md:flex">
                <div class="font-bold text-gray-500">me</div>
              </div>
            </td>
            <td class="py-3 space-x-1">
              <span v-if="user.role === 'admin'" class="bg-green-100 text-green-800 text-xs px-2 py-1 rounded">Admin</span>
              <span v-else class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">User</span>
            </td>
            <td class="py-3 hidden lg:table-cell">{{ new Date(user.UpdatedAt).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }) }}</td>
            <td class="py-3 hidden lg:table-cell">{{ new Date(user.CreatedAt).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }) }}</td>
            
            <td class="flex py-3 justify-start pr-3">
              <button
                @click="sendMessage(user)"
                class="w-1/2 h-1/2 bg-green-600 text-white rounded-md transition mx-1 p-1
                        hover:bg-green-500 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                <i class="pi pi-comments px-1"></i>
                <span class="hidden lg:inline"><span class="hidden xl:inline">Send</span> Message</span>
              </button>
              <button
                v-if="userAccoutStore.user.role === 'admin' "
                @click="toggleUpdateModal(user.ID,true)"
                :disabled="user.ID === userAccoutStore.user.id"
                class="w-1/2 h-1/2 bg-black text-white rounded-md transition mx-1 p-1
                        hover:bg-gray-800 cursor-pointer disabled:cursor-not-allowed hover:disabled:bg-black"
                >
                <i class="pi pi-pen-to-square px-1"></i>
                <span class="hidden lg:inline">Edit User</span>
              </button>
              <button 
                v-if="userAccoutStore.user.role === 'admin' "
                @click="toggleConfirmDeleteModal(user.ID, false)" 
                :disabled="user.ID === userAccoutStore.user.id"
                class="w-1/2 h-1/2 bg-red-800 text-white rounded-md transition mx-1 p-1 
                    hover:bg-red-900 cursor-pointer  disabled:cursor-not-allowed hover:disabled:bg-red-800"
                title="Delete User"
              >
                <i class="pi pi-trash px-1"></i>
                <span class="hidden lg:inline">Delete<span class="hidden xl:inline"> User</span></span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  <UserUpdateModal v-if="isShowUpdateModal" :toggle-modal="toggleUpdateModal" :user-id="selectedUserId" :is-edit-mode="isEditMode"></UserUpdateModal>
  <ConfirmDeleteModal v-if="isShowConfirmDeleteModal" :toggle-modal="toggleConfirmDeleteModal" :user-id="selectedUserId" ></ConfirmDeleteModal>
</template>


