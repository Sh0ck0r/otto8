import {
    ClientLoaderFunctionArgs,
    redirect,
    useLoaderData,
    useNavigate,
} from "@remix-run/react";
import { useCallback } from "react";
import { $path } from "remix-routes";
import { z } from "zod";

import { AgentService } from "~/lib/service/api/agentService";
import { RouteService } from "~/lib/service/routeQueryParams";
import { noop } from "~/lib/utils";

import { Agent } from "~/components/agent";
import { AgentProvider } from "~/components/agent/AgentContext";
import { Chat, ChatProvider } from "~/components/chat";
import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup,
} from "~/components/ui/resizable";

export type SearchParams = z.infer<
    (typeof RouteService.schemas)["/agents/:agent"]
>;

export const clientLoader = async ({
    params,
    request,
}: ClientLoaderFunctionArgs) => {
    const url = new URL(request.url);

    const { agent: agentId } = RouteService.getPathParams(
        "/agents/:agent",
        params
    );

    const { threadId, from } =
        RouteService.getQueryParams("/agents/:agent", url.search) ?? {};

    if (!agentId) {
        throw redirect("/agents");
    }

    // preload the agent
    const agent = await AgentService.getAgentById(agentId).catch(noop);

    if (!agent) {
        throw redirect("/agents");
    }
    return { agent, threadId, from };
};

export default function ChatAgent() {
    const { agent, threadId } = useLoaderData<typeof clientLoader>();
    const navigate = useNavigate();

    const updateThreadId = useCallback(
        (newThreadId?: Nullish<string>) => {
            navigate(
                $path(
                    "/agents/:agent",
                    { agent: agent.id },
                    newThreadId ? { threadId: newThreadId } : undefined
                )
            );
        },
        [agent, navigate]
    );

    return (
        <div className="h-full flex flex-col overflow-hidden relative">
            <ChatProvider
                id={agent.id}
                threadId={threadId}
                onCreateThreadId={updateThreadId}
            >
                <AgentProvider agent={agent}>
                    <ResizablePanelGroup
                        direction="horizontal"
                        className="flex-auto"
                    >
                        <ResizablePanel className="">
                            <Agent agent={agent} onRefresh={updateThreadId} />
                        </ResizablePanel>
                        <ResizableHandle withHandle />
                        <ResizablePanel>
                            <Chat />
                        </ResizablePanel>
                    </ResizablePanelGroup>
                </AgentProvider>
            </ChatProvider>
        </div>
    );
}
