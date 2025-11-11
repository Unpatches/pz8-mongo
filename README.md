## Имя: Дорджиев Виктор
## Группа: ЭФМО-02-25
# Проект pz8-mongo

Задачи проекта
1)	Понять базовые принципы документной БД MongoDB (документ, коллекция, BSON, _id:ObjectID).
2)	Научиться подключаться к MongoDB из Go с использованием официального драйвера.
3)	Создать коллекцию, индексы и реализовать CRUD для одной сущности (например, notes).
4)	Отработать фильтрацию, пагинацию, обновления (в т.ч. частичные), удаление и обработку ошибок.


---

## Установка и запуск

(Необходимы предустановленные Go и Git)

Клонировать репозиторий:

```
git clone https://github.com/Unpatches/pz8-mongo
cd pz8-mongo
```

Команда запуска сервера:

```
go run ./cmd/server
```


------

## Структура проекта

```plaintext
pz8-mongo/                     
├── cmd/                  
│   └── server/             
│       └── main.go       
├── internal/              
│   └── db/               
│   │   └── mongo.go
│   └── notes/               
│       └── handler.go
│       └── model.go
│       └── repo.go          
├── go.mod                
└── go.sum       
```

## Отчёт о проделанной работе

<img width="480" height="501" alt="image" src="https://github.com/user-attachments/assets/5ad99c4d-032a-446c-adb8-acd653ce42eb" />


<img width="757" height="686" alt="image" src="https://github.com/user-attachments/assets/223313b8-5bab-4c69-a500-849b0e6fc29c" />


<img width="739" height="319" alt="image" src="https://github.com/user-attachments/assets/0cd07951-6af8-4665-8cbc-704cea1275e9" />


<img width="494" height="467" alt="image" src="https://github.com/user-attachments/assets/4cf9482d-2b85-400a-95bd-618442ec5cc1" />


<img width="498" height="375" alt="image" src="https://github.com/user-attachments/assets/aedd1bc1-be95-4465-976d-a25dafd93fde" />





