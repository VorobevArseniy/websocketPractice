import React, { useState, useEffect } from "react";
import { View, Text, StyleSheet } from "react-native";
import {
	Camera,
	useCameraDevices,
	useCodeScanner,
} from "react-native-vision-camera";

interface QRScannerProps {
	onScan: (data: string | null) => void;
}

export default function QRScannerScreen({ onScan }: QRScannerProps) {
	const [hasPermission, setHasPermission] = useState<boolean | null>(null);
	const devices = useCameraDevices();
	const device = devices.find((d) => d.position === "back");

	useEffect(() => {
		(async () => {
			const status = await Camera.requestCameraPermission();
			setHasPermission(status === "granted");
		})();
	}, []);

	const codeScanner = useCodeScanner({
		codeTypes: ["qr"],
		onCodeScanned: (codes) => {
			const scannedValue = codes[0]?.value ?? null;
			console.log("QR Code:", scannedValue);
			onScan(scannedValue); // Передаём данные в родительский компонент
		},
	});

	if (hasPermission === null) {
		return <Text>Запрос разрешений...</Text>;
	}

	if (!hasPermission) {
		return <Text>Доступ к камере запрещен</Text>;
	}

	if (!device) {
		return <Text>Камера не найдена</Text>;
	}

	return (
		<View style={styles.container}>
			<Camera
				style={StyleSheet.absoluteFill}
				device={device}
				isActive={true}
				codeScanner={codeScanner}
			/>
			<View style={styles.overlay}>
				<Text style={styles.text}>Наведи камеру на QR-код</Text>
			</View>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
	},
	overlay: {
		position: "absolute",
		bottom: 20,
		left: 0,
		right: 0,
		alignItems: "center",
	},
	text: {
		fontSize: 18,
		color: "white",
		backgroundColor: "rgba(0, 0, 0, 0.7)",
		padding: 10,
		borderRadius: 5,
	},
});
