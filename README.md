## Логика работы программы:

#### 1. Чтение входных данных

Приложение начинает с чтения аргументов командной строки, чтобы определить имя входного файла.
После получения имени файла вызывается функция ParseInput из пакета util, которая читает файл и разбирает его содержимое. Входные данные содержат информацию о количестве муравьев, комнатах, туннелях и других деталях муравейника.

#### 2. Создание муравейника

После чтения входных данных вызывается метод NewAntFarm и создается муравейник (AntFarm)

#### 3. Нахождение всех путей

После создания фермы вызывается метод FindAllPaths у созданного муравейника (AntFarm). Этот метод использует алгоритм поиска в глубину (BFS), чтобы найти все возможные пути от начальной до конечной комнаты в муравейнике.

#### 4. Группировка путей

Найденные пути затем группируются с помощью функции FindNonIntersectingPathGroups, которая формирует группы путей, которые не пересекаются друг с другом. Это важно для эффективного распределения муравьев по путям.

#### 5. Выбор лучшей группы
Далее выбирается оптимальная группа путей с использованием функции ChooseBestGroup. Она определяет группу с наименьшим временем пересечения для всех муравьев.

#### 6. Распределение муравьев
После выбора оптимальной группы вызывается функция DistributeAnts, которая распределяет муравьев с учетом минимизации времени прохода. Это позволяет эффективно использовать все пути и максимизировать производительность.

#### 7. Симуляция движения муравьев
Наконец, вызывается функция SimulateAntsMovement, которая моделирует движение муравьев по муравейнику. Она выводит последовательность шагов, показывающую, как каждый муравей перемещается по своему пути к конечной комнате.
В итоге приложение решает задачу направления муравьев от начальной до конечной комнаты в муравейнике, учитывая ограничения на количество муравьев и длину путей.