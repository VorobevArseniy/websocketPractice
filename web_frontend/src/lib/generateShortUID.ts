import { customAlphabet } from "nanoid";

export function generateUID(size: number): string {
	const nanoid = customAlphabet(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
		size,
	);

	return nanoid();
}
