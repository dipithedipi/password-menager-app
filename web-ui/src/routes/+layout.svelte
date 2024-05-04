<script>
	import '../app.pcss';
 	import { DarkMode } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';

	// check if the route is different from /login or /register
	let currentUrl = $page.url.pathname;
	let contentPage = true
	if (currentUrl !== '/login' && currentUrl !== '/register') {
		contentPage = false 
	}
</script>

<!-- <DarkMode /> -->
<div class="container mx-auto py-2 overflow-hidden">
	<Navbar></Navbar>
	{#if !contentPage}
		<div class="flex pt-2 space-x-2 h-full">
			<Sidebar></Sidebar>
			<slot></slot>
	  </div>
	{:else}
		<slot></slot>
	{/if}
</div>
