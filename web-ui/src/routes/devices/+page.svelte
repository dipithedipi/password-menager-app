<script>
	import PhoneDevice from '$lib/components/PhoneDevice.svelte';
</script>

<div class="h-screen pt-6 bg-gray-700">
    <PhoneDevice activeSession=true></PhoneDevice>
    <PhoneDevice activeSession=false></PhoneDevice>
    <PhoneDevice activeSession=false></PhoneDevice>
</div>