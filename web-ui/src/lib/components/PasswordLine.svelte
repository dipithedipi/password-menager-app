<script lang="ts">
	export let modalTarget: string;
</script>

<tr
	class="border-b hover:bg-gray-600 dark:border-gray-700"
	data-modal-target={modalTarget}
	data-modal-toggle={modalTarget}
>
	<td class="h-10 w-10 px-4 py-4">
		<img src="/logo.jpg" alt="" />
	</td>
	<th scope="row" class="whitespace-nowrap px-4 py-4 font-medium text-gray-900 dark:text-white">
		test@mail.com
		<div class="text-slate-400 sm:hidden">nvidia.com</div>
	</th>
	<td class="hidden px-4 py-4 sm:flex">nvidia.com</td>
	<td class="px-4 py-4">3 days ago</td>
	<td class="items-center justify-end px-4 py-4">
		<button
			id="apple-imac-27-dropdown-button"
			data-dropdown-toggle="apple-imac-27-dropdown"
			class="inline-flex items-center rounded-lg p-0.5 text-center text-sm font-medium text-gray-500 hover:text-gray-800 focus:outline-none dark:text-gray-400 dark:hover:text-gray-100"
			type="button"
		>
			<svg
				class="h-5 w-5"
				aria-hidden="true"
				fill="currentColor"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z"
				/>
			</svg>
		</button>
		<div
			id="apple-imac-27-dropdown"
			class="z-10 hidden w-44 divide-y divide-gray-100 rounded bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
		>
			<ul
				class="py-1 text-sm text-gray-700 dark:text-gray-200"
				aria-labelledby="apple-imac-27-dropdown-button"
			>
				<li>
					<a
						href="#"
						class="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white"
						>Show</a
					>
				</li>
				<li>
					<a
						href="#"
						class="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white"
						>Edit</a
					>
				</li>
			</ul>
			<div class="py-1">
				<a
					href="#"
					class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-600 dark:hover:text-white"
					>Delete</a
				>
			</div>
		</div>
	</td>
</tr>
