<script>
	import '../app.pcss';
 	import { DarkMode } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Navbar from '$lib/components/Navbar.svelte';

	$: navbarVisible = $page.url.pathname !== '/login' && $page.url.pathname !== '/register' && $page.url.pathname !== '/';
</script>

<!-- <DarkMode /> -->
<div class="bg-slate-600">
	<div>
		{#if navbarVisible}
			<div>
				<Sidebar><slot></slot></Sidebar>
			</div>
		{:else}
			<div class="h-full">
				<Navbar/>
				<slot></slot>
			</div>
		{/if}
	</div>
</div>
