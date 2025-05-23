<script lang="ts">
	import { Thread } from '$lib/services/chat/thread.svelte';
	import type { Messages, Project, ProjectCredential } from '$lib/services';
	import Message from '$lib/components/messages/Message.svelte';
	import { X } from 'lucide-svelte';
	import { responsive } from '$lib/stores';
	import { twMerge } from 'tailwind-merge';

	interface Props {
		toolID: string;
		onClose?: () => void | Promise<void>;
		project: Project;
		local?: boolean;
		credential?: ProjectCredential;
	}

	let { toolID, onClose, project, local, credential }: Props = $props();
	let authMessages = $state<Messages>();
	let thread = $state<Thread>();
	let authDialog: HTMLDialogElement | undefined = $state();

	export function show() {
		authCancel(true);
		auth(toolID);
		authDialog?.showModal();
	}

	function auth(toolID: string) {
		const t = new Thread(project, {
			authenticate: {
				tools: [toolID],
				local
			},
			onError: () => {
				// ignore the error. This is so it doesn't get globally printed
			},
			onClose: () => {
				authCancel();
				// false means don't reconnect
				return false;
			}
		});
		t.onMessages = (messages) => {
			authMessages = messages;
		};
		thread = t;
	}

	function authCancel(skipClose?: boolean) {
		thread?.abort();
		thread?.close();
		// only clear message if nothing failed
		if (!(authMessages?.messages ?? []).find((msg) => msg.icon === 'Error')) {
			authMessages = undefined;
		}
		thread = undefined;
		if (!skipClose) {
			authDialog?.close();
			onClose?.();
		}
	}
</script>

{#if local}
	{@render content()}
{:else}
	<dialog
		bind:this={authDialog}
		class:mobile-screen-dialog={responsive.isMobile}
		class="md:max-w-sm"
	>
		{@render content()}
	</dialog>
{/if}

{#snippet content()}
	{#if credential}
		<div class="flex flex-col">
			<h4
				class={twMerge(
					'default-dialog-title text-base',
					!local && 'py-2 md:pr-2 md:pl-5',
					local && 'mb-4'
				)}
				class:default-dialog-mobile-title={responsive.isMobile}
			>
				<span class="flex items-center gap-2">
					<img
						src={credential.icon}
						class="size-6 rounded-md bg-white p-1"
						alt="credential {credential.toolName} icon"
					/>
					{credential.toolName}
				</span>
				<button
					class={twMerge('icon-button', local && 'translate-x-2')}
					class:mobile-header-button={responsive.isMobile}
					onclick={() => authCancel()}
				>
					<X class="size-5" />
				</button>
			</h4>
			{#if authMessages}
				<div class={twMerge('flex flex-col gap-5 p-5', !local && 'md:m-5 md:mt-0')}>
					{#each authMessages.messages as msg}
						<Message {msg} {project} clearable onSendCredentialsCancel={() => authCancel()} />
					{/each}
				</div>
			{/if}
		</div>
	{/if}
{/snippet}
