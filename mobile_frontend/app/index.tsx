import React, { useState } from "react";
import { View, Text, Button, StyleSheet } from "react-native";
import QRScanner from "../components/QRScanner";

export default function App() {
	const [qrData, setQrData] = useState<string | null>(null);

	return (
		<View style={styles.container}>
			{qrData ? (
				<>
					<Text style={styles.text}>Отсканировано: {qrData}</Text>
					<Button title="Сканировать снова" onPress={() => setQrData(null)} />
				</>
			) : (
				<QRScanner onScan={setQrData} />
			)}
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		alignItems: "center",
	},
	text: {
		fontSize: 18,
		marginBottom: 10,
	},
});
