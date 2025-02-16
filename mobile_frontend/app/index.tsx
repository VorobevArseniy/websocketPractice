import React, { useState } from "react";
import { View, TextInput, Button, Alert, StyleSheet } from "react-native";

const App = () => {
  const [name, setName] = useState("");
  const [sessionID, setSessionID] = useState("");

  const handleSubmit = async () => {
    if (!name || !sessionID) {
      Alert.alert("Ошибка", "Заполните все поля");
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/adduser?sessionID=${sessionID}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ id: name }),
        },
      );

      if (!response.ok) {
        throw new Error(`Ошибка: ${response.statusText}`);
      }

      Alert.alert("Успех", "Пользователь добавлен!");
    } catch (error) {
      console.error("Ошибка при отправке запроса:", error);
      Alert.alert("Ошибка", "Не удалось отправить данные");
    }
  };

  return (
    <View style={styles.container}>
      <TextInput
        style={styles.input}
        placeholder="Введите ваше имя"
        value={name}
        onChangeText={setName}
      />
      <TextInput
        style={styles.input}
        placeholder="Введите ID сессии"
        value={sessionID}
        onChangeText={setSessionID}
      />
      <Button title="Отправить" onPress={handleSubmit} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    padding: 20,
    backgroundColor: "#f4f4f4",
  },
  input: {
    height: 50,
    borderColor: "#ccc",
    borderWidth: 1,
    marginBottom: 20,
    paddingHorizontal: 10,
    backgroundColor: "#fff",
    borderRadius: 5,
  },
});

export default App;
