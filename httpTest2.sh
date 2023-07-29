#!/bin/bash

# Задаем количество повторений (1000 в данном случае)
number_of_requests=1000

# Выполняем цикл для отправки POST-запросов
for ((i=1; i<=number_of_requests; i++)); do
  echo "Отправка запроса $i..."
  curl --location --request POST 'http://localhost:8080/rest/counter/sub/1'
  echo "Ответ: $response"
done