<script lang="ts">
	import { Plus, Minus } from 'lucide-svelte';
	import { autoHeight } from '$lib/actions/textarea';

	interface Props {
		params: { key: string; value: string }[];
		input?: boolean;
		autofocus?: boolean;
	}

	let { params = $bindable([]), input }: Props = $props();
</script>

<div class="bg-surface1 flex flex-col gap-4 rounded-lg p-5">
	<div class="flex">
		{#if !input}
			<h4 class="flex-1 text-xl font-semibold">Arguments</h4>
		{/if}
		{#if !input}
			<button onclick={() => params.push({ key: '', value: '' })}>
				<Plus class="h-5 w-5" />
			</button>
		{/if}
	</div>
	{#if params.length !== 0}
		<table class="w-full table-auto text-left">
			<thead>
				<tr>
					<th class="w-1/4">Name</th>
					{#if input}
						<th class="w-full">Value</th>
					{:else}
						<th class="w-full">Description</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each params as param, i}
					<tr>
						<td class="pr-2">
							<input
								bind:value={param.key}
								readonly={input}
								placeholder="Enter name"
								class={[
									'focus:ring-blue me-1 w-full rounded-lg p-2 outline-none focus:ring-2',
									input ? 'bg-surface1' : 'bg-surface2'
								]}
							/>
						</td>
						<td class="flex items-center pr-2">
							<textarea
								use:autoHeight
								class="text-input resize-none"
								rows="1"
								bind:value={param.value}
							></textarea>
						</td>
						<td>
							{#if !input}
								<button onclick={() => params.splice(i, 1)}>
									<Minus class="h-5 w-5" />
								</button>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>
