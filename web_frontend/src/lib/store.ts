import { createSignal } from "solid-js";

export interface formDataInterface {
	name: string;
	subject: string;
	cabinet: string;
}

export const [formData, setFormData] = createSignal<formDataInterface>({
	name: "",
	subject: "",
	cabinet: "",
});

export const [sessionID, setSessionID] = createSignal<string>("");
