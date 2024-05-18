<script>
	import EventLine from '$lib/components/EventLine.svelte';
</script>

<div class="pt-3 pb-3 h-full bg-gray-700">
	<div date-rangepicker class="flex items-center justify-center pt-6">
		<div class="ml-4 sm:ml-0 relative">
			<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
				<svg
					class="h-4 w-4 text-gray-500 dark:text-gray-400"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="currentColor"
					viewBox="0 0 20 20"
				>
					<path
						d="M20 4a2 2 0 0 0-2-2h-2V1a1 1 0 0 0-2 0v1h-3V1a1 1 0 0 0-2 0v1H6V1a1 1 0 0 0-2 0v1H2a2 2 0 0 0-2 2v2h20V4ZM0 18a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8H0v10Zm5-8h10a1 1 0 0 1 0 2H5a1 1 0 0 1 0-2Z"
					/>
				</svg>
			</div>
			<input
				name="start"
				type="text"
				class="block w-full rounded-lg border-2 border-gray-300 bg-gray-50 p-2.5 ps-10 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
				placeholder="Select date start"
			/>
		</div>
		<span class="mx-3 sm:mx-4 text-gray-400">to</span>
		<div class="mr-4 sm:mr-0 relative">
			<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
				<svg
					class="h-4 w-4 text-gray-500 dark:text-gray-400"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="currentColor"
					viewBox="0 0 20 20"
				>
					<path
						d="M20 4a2 2 0 0 0-2-2h-2V1a1 1 0 0 0-2 0v1h-3V1a1 1 0 0 0-2 0v1H6V1a1 1 0 0 0-2 0v1H2a2 2 0 0 0-2 2v2h20V4ZM0 18a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8H0v10Zm5-8h10a1 1 0 0 1 0 2H5a1 1 0 0 1 0-2Z"
					/>
				</svg>
			</div>
			<input
				name="end"
				type="text"
				class="block w-full rounded-lg border-2 border-gray-300 bg-gray-50 p-2.5 ps-10 text-sm text-gray-900 focus:border-gray-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-white dark:placeholder-gray-400 dark:focus:border-gray-500 dark:focus:ring-blue-500"
				placeholder="Select date end"
			/>
		</div>
	</div>
	<div class="w-3/4 flex flex-col items-center justify-center pt-4 mx-auto">
		<EventLine/>
	</div>
</div>
