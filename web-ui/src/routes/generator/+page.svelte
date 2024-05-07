<div class="h-screen justify-center bg-gray-700">
	<div class="pt-8 sm:pt-7 mx-2">
		<div
			class="mx-auto flex max-w-lg items-center rounded-lg border border-gray-700 bg-gray-800 p-3"
		>
			<div class="relative w-full">
				<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
					<svg
						class="h-4 w-4 text-gray-500 dark:text-gray-400"
						aria-hidden="true"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 21 21"
					>
						<path
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M11.15 5.6h.01m3.337 1.913h.01m-6.979 0h.01M5.541 11h.01M15 15h2.706a1.957 1.957 0 0 0 1.883-1.325A9 9 0 1 0 2.043 11.89 9.1 9.1 0 0 0 7.2 19.1a8.62 8.62 0 0 0 3.769.9A2.013 2.013 0 0 0 13 18v-.857A2.034 2.034 0 0 1 15 15Z"
						/>
					</svg>
				</div>
				<input
					type="text"
					id="voice-search"
					class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 ps-10 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
					placeholder="Random Password"
					required
				/>
			</div>
			<button
				type="submit"
				class="ms-2 inline-flex items-center rounded-lg border border-blue-700 bg-blue-700 px-3 py-2.5 text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
			>
				<svg
					class="me-2 h-4 w-4"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 20 20"
				>
					<path
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
					/>
				</svg>New
			</button>
			<button
				type="submit"
				class="ms-2 inline-flex items-center rounded-lg border border-blue-700 bg-blue-700 px-3 py-2.5 text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
			>
				<svg
					class="me-2 h-4 w-4"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 20 20"
				>
					<path
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
					/>
				</svg>Copy
			</button>
		</div>
	</div>

	<div class="mx-2 md:mx-auto max-w-lg pt-2">
		<div class="relative max-w-lg rounded-lg border-gray-700 bg-gray-800 p-3">
            <div
				class="flex items-center font-medium rounded-lg border border-gray-300 bg-gray-50 p-2.5 py-3 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
			>
				Lenght:

				<div class="absolute right-6 flex max-w-[8rem] items-center">
					<button
						type="button"
						id="decrement-button"
						data-input-counter-decrement="quantity-input"
						class="h-10 rounded-s-lg border border-gray-300 bg-gray-100 p-2 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:focus:ring-gray-700"
					>
						<svg
							class="h-3 w-3 text-gray-900 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 18 2"
						>
							<path
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M1 1h16"
							/>
						</svg>
					</button>
					<input
						type="text"
						id="quantity-input"
						data-input-counter
						aria-describedby="helper-text-explanation"
						class="block h-10 w-14 border-x-0 border-gray-300 bg-gray-50 py-2.5 text-center text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
						placeholder="32"
						required
					/>
					<button
						type="button"
						id="increment-button"
						data-input-counter-increment="quantity-input"
						class="h-10 rounded-e-lg border border-gray-300 bg-gray-100 p-2 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:focus:ring-gray-700"
					>
						<svg
							class="h-3 w-3 text-gray-900 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 18 18"
						>
							<path
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 1v16M1 9h16"
							/>
						</svg>
					</button>
				</div>
			</div>

            <div class="mt-3 flex items-center rounded border border-gray-200 ps-4 dark:border-gray-700">
                <input
					checked
					id="bordered-checkbox-4"
					type="checkbox"
					value=""
					name="bordered-checkbox"
					class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
				/>
				<label
					for="bordered-checkbox-4"
					class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
					>Lowercase</label
				>
			</div>
			<div class="my-3 flex items-center rounded border border-gray-200 ps-4 dark:border-gray-700">
				<input
					id="bordered-checkbox-3"
					type="checkbox"
					value=""
					name="bordered-checkbox"
					class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
				/>
				<label
					for="bordered-checkbox-3"
					class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
					>Uppercase</label
				>
			</div>
			<div class="my-3 flex items-center rounded border border-gray-200 ps-4 dark:border-gray-700">
				<input
					id="bordered-checkbox-1"
					type="checkbox"
					value=""
					name="bordered-checkbox"
					class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
				/>
				<label
					for="bordered-checkbox-1"
					class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
					>Numbers</label
				>
			</div>
			<div class="my-3 flex items-center rounded border border-gray-200 ps-4 dark:border-gray-700">
				<input
					id="bordered-checkbox-2"
					type="checkbox"
					value=""
					name="bordered-checkbox"
					class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
				/>
				<label
					for="bordered-checkbox-2"
					class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
					>Special character</label
				>
			</div>
            
		</div>
        <div class="flex pt-3 items-center rounded-lg border-gray-700 bg-gray-800 p-3 mt-2"> 
            <div class="ms-2 w-full text-sm font-medium text-gray-900 dark:text-gray-300">
                Security level:
            </div>
            <span class="bg-green-100 text-green-800 text-md font-medium me-2 px-2.5 py-0.5 rounded dark:bg-green-900 dark:text-green-300">Strong</span>
        </div>
	</div>
   
</div>
