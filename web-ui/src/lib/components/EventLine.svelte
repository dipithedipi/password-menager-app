<script lang="ts">
	export let type: string;
	export let id: string;
	export let title: string;
	export let description: string;
	export let ip: string;
	export let time: string;

	let classEventMap: any = {
		"info": "bg-blue-50 border-blue-300 text-blue-800 dark:bg-blue-800 dark:border-blue-800 dark:text-blue-400",
		"danger": "bg-red-50 border-red-300 text-red-800 dark:bg-blue-800 dark:border-red-800 dark:text-red-400",
		"success": "bg-green-50 border-green-300 text-green-800 dark:bg-blue-800 dark:border-green-800 dark:text-green-400",
		"warning": "bg-yellow-50 border-yellow-300 text-yellow-800 dark:bg-blue-800 dark:border-yellow-800 dark:text-yellow-300",
		"default": "bg-gray-50 border-gray-300 text-gray-800 dark:bg-blue-800 dark:border-gray-800 dark:text-gray-400"
	};

	let classButtonMap: any = {
		"info": "bg-blue-800 text-white hover:bg-blue-900 focus:ring-blue-200 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 dark:border-blue-800",
		"danger": "bg-red-800 text-white hover:bg-red-900 focus:ring-red-300 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-800 dark:border-red-800",
		"success": "bg-green-800 text-white hover:bg-green-900 focus:ring-green-300 dark:bg-green-600 dark:hover:bg-green-700 dark:focus:ring-green-800 dark:border-green-800",
		"warning": "bg-yellow-800 text-white hover:bg-yellow-900 focus:ring-yellow-300 dark:bg-yellow-300 dark:hover:bg-yellow-400 dark:focus:ring-yellow-800 dark:border-yellow-800",
		"default": "bg-gray-700 text-white hover:bg-gray-800 focus:ring-gray-300 dark:bg-gray-600 dark:hover:bg-gray-500 dark:focus:ring-gray-800 dark:border-gray-800"
	}

	function dismissAlert(idElm: string) {
		const alertEvent: any = document.getElementById(`${idElm}-event`);
		if (!alertEvent) {
			alert("event not found")
			return
		}
		alertEvent.style.display = "none";
	}

	function showMore(ipEvent: string) {
		alert("IP: " + ipEvent + "\nTime: " + time + "\nDescription: " + description)
	}
</script>


<div
	id="{id}-event"
	class={`mb-4 w-11/12 rounded-lg border-2 ${classEventMap[type]} p-4 text-blue-800 dark:border-blue-800 dark:bg-gray-800 dark:text-blue-400`}
	role={type}
>
	<div class="flex items-center">
		<svg
			class="me-2 h-4 w-4 flex-shrink-0"
			aria-hidden="true"
			xmlns="http://www.w3.org/2000/svg"
			fill="currentColor"
			viewBox="0 0 20 20"
		>
			<path
				d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM9.5 4a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM12 15H8a1 1 0 0 1 0-2h1v-3H8a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1v4h1a1 1 0 0 1 0 2Z"
			/>
		</svg>
		<span class="sr-only">{type}</span>
		<h3 class="text-lg font-medium">{title} ({type} alert)</h3>
	</div>
	<div class="mb-4 mt-2 text-sm">
		{description}
	</div>
	<div class="flex space-x-3">
		<button
			type="button"
			on:click={() => showMore(ip)}
			class={`flex rounded-lg border ${classButtonMap[type]} px-3 py-1.5 text-center text-xs font-medium`}
		>
			<svg
				class="me-2 h-3 w-3"
				aria-hidden="true"
				xmlns="http://www.w3.org/2000/svg"
				fill="currentColor"
				viewBox="0 0 20 14"
			>
				<path
					d="M10 0C4.612 0 0 5.336 0 7c0 1.742 3.546 7 10 7 6.454 0 10-5.258 10-7 0-1.664-4.612-7-10-7Zm0 10a3 3 0 1 1 0-6 3 3 0 0 1 0 6Z"
				/>
			</svg>
			View more
		</button>
		<button
			on:click={() => dismissAlert(id)}
			type="button"
			class={`rounded-lg border ${classButtonMap[type]} px-3 py-1.5 text-center text-xs font-medium `}
			data-dismiss-target="#alert-additional-content-1"
			aria-label="Close"
		>
			Dismiss
		</button>
	</div>
</div>
