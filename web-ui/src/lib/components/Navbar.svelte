<script>
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';

    // @ts-ignore
    function handleGetStartedClick(event) {
        event.preventDefault();
        goto('/login');
    }
   
    $: classAttr = "flex md:order-2 space-x-3 md:space-x-0 rtl:space-x-reverse " + ($page.url.pathname === "/" ? "" : "hidden");
</script>

<div class="fixed h-16 w-screen">
    <nav class="bg-white border-gray-200 dark:bg-gray-800">
        <div class="flex flex-wrap items-center justify-between mx-auto p-4">
        <a href="/" class="flex items-center space-x-3 rtl:space-x-reverse">
            <img src="/logo.jpg" class="h-9 rounded-full" alt="Logo" />
            <span class="self-center text-2xl font-semibold whitespace-nowrap dark:text-white">PassGuard</span>
        </a>
        <div class={classAttr}>
            <button type="button" on:click={handleGetStartedClick} class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
            Get started
            </button>
        </div>
        </div>
    </nav>
</div>