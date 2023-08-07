"use client";

export type Message = {
  content: string;
  clientId: string;
  username: string;
  roomId: string;
  type: "recv" | "self";
};

export function ChatBody({ messages }: { messages: Array<Message> }) {
  return (
    <div>
      {messages.map((message: Message, index: number) => {
        if (message.type == "self") {
          return (
            <div
              className="flex flex-col mt-2 w-full text-right justify-end"
              key={index}
            >
              <div className="text-sm">{message.username}</div>
              <div>
                <div className="bg-green text-white px-4 py-1 rounded-md inline-block mt-1">
                  {message.content}
                </div>
              </div>
            </div>
          );
        } else {
          return (
            <div className="mt-2" key={index}>
              <div className="text-sm">{message.username}</div>
              <div>
                <div className="bg-grey text-dark-secondary px-4 py-1 rounded-md inline-block mt-1">
                  {message.content}
                </div>
              </div>
            </div>
          );
        }
      })}
    </div>
  );
}
