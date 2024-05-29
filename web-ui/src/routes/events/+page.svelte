<script lang="ts">
	import { getEvents } from '$lib/logic/fetch';
	import EventLine from '$lib/components/EventLine.svelte';
	import { formatDate } from '$lib/logic/utils';
	import { onMount } from 'svelte';

	let events: any[] = [];
	let invalidDate: boolean = false;
	let notFound: boolean = false;

	let start: string = '';
	let end: string = '';

	async function handleEvents() {
		// reset
		events = [];
		notFound = false;

		start = (document.getElementById('start-date') as HTMLInputElement).value;
		end = (document.getElementById('end-date') as HTMLInputElement).value;	
		
		let startFormatted = formatDate(start) + "T00:00:00.000Z";
		let endFormatted = formatDate(end) + "T23:59:59.000Z";

		let {data, success} = await getEvents(startFormatted, endFormatted);
		if (!success) {
			notFound = true;
			return;
		}
		events = data.events;
		events.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
	}

	function formatStartDate(date: Date): string {
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are zero-based in JavaScript
		const day = String(date.getDate()).padStart(2, '0');

		return `${year}-${month}-${day}`;
	}

	function getTodayAndOneWeekAgoFormatted(): [string, string] {
		const today = new Date();
		const oneWeekAgo = new Date(today.getTime());
		oneWeekAgo.setDate(today.getDate() - 7);

		// Safely set the values only if the elements exist
        const startDateElement: any = document.getElementById('start-date');
        const endDateElement: any = document.getElementById('end-date');

        if (startDateElement && endDateElement) {
            startDateElement.value = formatStartDate(oneWeekAgo).split('-').reverse().join('/');
            endDateElement.value = formatStartDate(today).split('-').reverse().join('/');
        }

		return [
			formatStartDate(today) + "T23:59:59.000Z",
			formatStartDate(oneWeekAgo) + "T00:00:00.000Z"
		];
	}

	onMount(async () => {
		const [start, end] = getTodayAndOneWeekAgoFormatted();
		
		let {data, success} = await getEvents(end, start);
		if (!success) {
			notFound = true;
			return;
		}
		events = data.events;
		events.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
	})
</script>

<div class="{events.length < 8 ? "h-screen": "h-full"} pt-3 pb-3 bg-gray-700">
	<div date-rangepicker datepicker-format="dd/mm/yyyy" class="ml-2 flex items-center justify-center pt-6">
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
				id="start-date"
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
				id="end-date"
				class="block w-full rounded-lg border-2 border-gray-300 bg-gray-50 p-2.5 ps-10 text-sm text-gray-900 focus:border-gray-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-white dark:placeholder-gray-400 dark:focus:border-gray-500 dark:focus:ring-blue-500"
				placeholder="Select date end"
			/>
		</div>
		<div class="mr-2 md:mx-4">
			<button
				on:click={handleEvents}
				type="button" 
				class="h-10 flex items-center justify-center text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-primary-600 dark:hover:bg-primary-700 focus:outline-none dark:focus:ring-primary-800"
			>
			<svg class="h-4 w-4 sm:mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
					<path d="M18.031 16.6168L22.3137 20.8995L20.8995 22.3137L16.6168 18.031C15.0769 19.263 13.124 20 11 20C6.032 20 2 15.968 2 11C2 6.032 6.032 2 11 2C15.968 2 20 6.032 20 11C20 13.124 19.263 15.0769 18.031 16.6168ZM16.0247 15.8748C17.2475 14.6146 18 12.8956 18 11C18 7.1325 14.8675 4 11 4C7.1325 4 4 7.1325 4 11C4 14.8675 7.1325 18 11 18C12.8956 18 14.6146 17.2475 15.8748 16.0247L16.0247 15.8748Z"></path>
				</svg>
				<div class="hidden sm:block">
					Search
				</div>
			</button>
		</div>
	</div>
	<div class="w-3/4 flex flex-col items-center justify-center pt-4 mx-auto">
		{#if invalidDate}
			<p class="text-red-500">Please select a valid date range</p>
		{:else if notFound}
			<p class="text-gray-400">No events found</p>
		{:else}
			{#if events.length === 0}
				<p class="text-gray-400">No events found</p>
			{:else}
				{#each events as event, index}
					<EventLine id={index.toString()} type={event.type} title={event.title} description={event.description} time={event.createdAt} ip={event.ipAddress}/>
				{/each}
			{/if}
		{/if}
	</div>
</div>
