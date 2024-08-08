package lemin

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Ant struct {
	ID       int
	Position int      // Текущая позиция в определенном пути
	Path     []string // Путь по которому пойдет муравей
}

type PathInfo struct {
	Path           []string // Комнаты пути
	Ants           int      // Количество муравьев на этот путь
	CompletionTime int      // Длина пути
}

// PathGroup группирует пути с максимальной длиной и количеством муравьев, которые могут идти этим путем.
type PathGroup struct {
	Paths       [][]string // Список путей
	MaxLength   int        // Максимальная длина одного пути в группе
	MaxAnts     int        // Максимальное количество муравьев, которые могут одновременно следовать этим путям
	CurrentAnts []int      // Текущее количество муравьев на каждом пути
}

// Room представляет из себя структуру комнаты в муравьиной ферме
type Room struct {
	Name string
	X, Y int // В данном проекте не используются
}

// AntFarm представляет муравьиную ферму с картой комнат и туннелей.
type AntFarm struct {
	Rooms      map[string]Room     // Карта всех комнат
	Tunnels    map[string][]string // Список туннелей (связей между комнатами)
	Start, End string              // Начальная и конечная комнаты
}

// Вспомогательная структура для хранения информации о состоянии поиска
type searchState struct {
	room string   // Текущая комната
	path []string // Путь от начала до текущей комнаты
}

// distributeAntsOneGroup рассчитывает минимальное число шагов для того что бы все муравьи посланные по этой группе дошли до конца/ Нужна для выбора лучшей группы путей
func distributeAntsOneGroup(totalAnts, maxAnts, maxLength, numPaths int) int {
	if totalAnts <= maxAnts {
		return maxLength
	}

	excessAnts := totalAnts - maxAnts
	additionalTurns := (excessAnts + numPaths - 1) / numPaths
	return maxLength + additionalTurns
}

// DistributeAnts равномерно распределяет муравьев по путям в выбранной группе.
func DistributeAnts(totalAnts int, group *PathGroup) []Ant {
	// Создаем слайс для хранения муравьев, которым будут назначены пути.
	ants := make([]Ant, totalAnts)

	// Инициализируем массив для отслеживания количества муравьев на каждом пути.
	group.CurrentAnts = make([]int, len(group.Paths))

	// Используем структуру минимальной кучи для управления заполнением путей.
	type pathSlot struct {
		index int // Индекс пути
		cost  int // Стоимость - текущая задержка, если добавить еще одного муравья
	}
	pathHeap := make([]pathSlot, len(group.Paths))

	// Инициализируем кучу начальными значениями стоимостей (длиной каждого пути).
	for i, path := range group.Paths {
		pathHeap[i] = pathSlot{
			index: i,
			cost:  len(path),
		}
	}

	// Отслеживаем размер кучи.
	heapSize := len(pathHeap)

	// Функция для сортировки кучи по стоимости (сначала идут пути с минимальной стоимостью).
	heapify := func() {
		sort.Slice(pathHeap, func(i, j int) bool {
			return pathHeap[i].cost < pathHeap[j].cost
		})
	}

	// Назначаем муравьев на пути.
	for i := 0; i < totalAnts; i++ {
		// Всегда назначаем муравья на путь с наименьшей текущей стоимостью.
		heapify()
		minPath := pathHeap[0]

		// Назначаем муравья.
		ants[i] = Ant{
			ID:       i + 1,
			Path:     group.Paths[minPath.index],
			Position: 0,
		}

		// Обновляем количество муравьев на этом пути.
		group.CurrentAnts[minPath.index]++

		// Пересчитываем стоимость добавления еще одного муравья на этот путь.
		pathHeap[0].cost = len(group.Paths[minPath.index]) + group.CurrentAnts[minPath.index] - 1

		// Если размер кучи больше 1, возможно, потребуется переупорядочить кучу.
		if heapSize > 1 && pathHeap[0].cost > pathHeap[1].cost {
			heapify()
		}
	}

	// Выводим распределение муравьев по путям.
	// fmt.Printf("Distribution of ants per path: %v\n", group.CurrentAnts)
	return ants
}

// ChooseBestGroup выбирает лучшую группу для минимизации общего времени
func ChooseBestGroup(ants int, groups []PathGroup) (int, int, [][]string) {
	minSteps := math.MaxInt32
	bestGroupIndex := -1
	var bestPaths [][]string

	for i, group := range groups {
		steps := distributeAntsOneGroup(ants, group.MaxAnts, group.MaxLength, len(group.Paths))
		if steps < minSteps {
			minSteps = steps
			bestGroupIndex = i
			bestPaths = group.Paths
		}
	}

	return bestGroupIndex, minSteps - 1, bestPaths // Не считаем стартовую комнату как шаг
}

func (farm *AntFarm) FindAllPaths() [][]string {
	var allPaths [][]string

	// Инициализация стартового состояния поиска
	startState := searchState{room: farm.Start, path: []string{farm.Start}}

	// Инициализация мапы посещенных комнат с отметкой стартовой комнаты как посещенной
	visitedStart := make(map[string]bool)
	visitedStart[farm.Start] = true

	var findPaths func(state searchState, visited map[string]bool)
	findPaths = func(state searchState, visited map[string]bool) {
		if state.room == farm.End {
			// Добавление найденного пути к листу allPaths
			allPaths = append(allPaths, state.path)
			return
		}

		// Цикл по всем соседним комнатам текущей комнаты
		for _, nextRoom := range farm.Tunnels[state.room] {
			if !visited[nextRoom] {
				// Создание копии текущего состояния посещенных комнат для следующего пути
				visitedCopy := make(map[string]bool)
				for key, value := range visited {
					visitedCopy[key] = value
				}
				visitedCopy[nextRoom] = true // Отметка о визите nextRoom

				// Создание нового пути с добавлением nextRoom
				newPath := make([]string, len(state.path)+1)
				copy(newPath, state.path)
				newPath[len(newPath)-1] = nextRoom

				// Продолжение поиска со следующей комнаты
				findPaths(searchState{room: nextRoom, path: newPath}, visitedCopy)
			}
		}
	}

	// Начинает поиск со стартовой комнаты используя мап со стартовой комнатой указанной как посещенная (DFS)
	findPaths(startState, visitedStart)

	return allPaths
}

// NewAntFarm создает новую AntFarm, из комнат и тунелей
func NewAntFarm(startRoom, endRoom string, rooms map[string]Room, tunnels []string) *AntFarm {
	farm := &AntFarm{
		Rooms:   rooms,
		Tunnels: make(map[string][]string),
		Start:   startRoom,
		End:     endRoom,
	}

	for _, tunnel := range tunnels {
		parts := strings.Split(tunnel, "-")
		if len(parts) == 2 {
			farm.Tunnels[parts[0]] = append(farm.Tunnels[parts[0]], parts[1])
			// For an undirected graph, add the reverse direction as well.
			farm.Tunnels[parts[1]] = append(farm.Tunnels[parts[1]], parts[0])
		}
	}

	return farm
}

// Проверяет, пересекаются ли два пути, исключая стартовую и конечную точки

func pathsIntersect(path1, path2 []string) bool {
	// Игнорируем стартовую и конечную точки в проверке пересечений
	set := make(map[string]bool)
	for _, node := range path1[1 : len(path1)-1] {
		set[node] = true
	}
	for _, node := range path2[1 : len(path2)-1] {
		if set[node] {
			return true
		}
	}
	// fmt.Println("Paths do not intersect")
	// fmt.Println(path1)
	// fmt.Println(path2)
	return false
}

// Возвращает все группы путей, которые не пересекаются друг с другом
func FindNonIntersectingPathGroups(paths [][]string) []PathGroup {
	var groups []PathGroup
	n := len(paths)

	if n == 0 {
		return groups // Если путей нет, возвращаем пустой список
	}

	// Используем функцию для создания всех возможных групп
	var findGroups func(currentIndex int, currentGroup []int)
	findGroups = func(currentIndex int, currentGroup []int) {
		// Если мы дошли до конца списка путей, создаем группу из текущих индексов
		if currentIndex == n {
			if len(currentGroup) == 0 {
				return
			}
			// Создаем новую группу из индексов currentGroup
			newGroup := make([][]string, len(currentGroup))
			maxLength := 0
			for i, idx := range currentGroup {
				newGroup[i] = paths[idx]
				if len(paths[idx]) > maxLength {
					maxLength = len(paths[idx])
				}
			}

			// Вычисляем максимальное количество муравьев для группы
			maxAnts := 0
			for _, path := range newGroup {
				maxAnts += maxLength - len(path) + 1
			}

			groups = append(groups, PathGroup{
				Paths: newGroup,
				//				PathIndexes: currentGroup,
				MaxLength: maxLength,
				MaxAnts:   maxAnts,
			})
			return
		}

		// Проверяем, можем ли мы добавить текущий путь в группу без конфликтов
		conflict := false
		for _, idx := range currentGroup {
			if pathsIntersect(paths[currentIndex], paths[idx]) {
				conflict = true
				break
			}
		}

		// Пробуем добавить путь в группу, если нет конфликтов
		if !conflict {
			findGroups(currentIndex+1, append(currentGroup, currentIndex))
		}

		// Пробуем пропустить текущий путь и продолжить поиск
		findGroups(currentIndex+1, currentGroup)
	}

	// Запускаем рекурсивное формирование групп начиная с пустой группы
	findGroups(0, []int{})

	// Сортируем группы по максимальной длине пути в каждой группе от меньшей к большей
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].MaxLength < groups[j].MaxLength
	})

	return groups
}

func SimulateAntsMovement(ants []Ant, startRoom, endRoom string) {
	fmt.Println("Simulation starts")
	turn := 0

	// Мапы для отслеживания занятости комнат и туннелей
	occupiedRooms := make(map[string]bool)
	occupiedTunnels := make(map[string]bool)

	// Флаг для проверки, все ли муравьи завершили движение
	allAntsFinished := false

	for !allAntsFinished {
		allAntsFinished = true
		moves := make([]string, 0)

		// Сбрасываем занятость туннелей на каждом шагу, т.к. муравей уже прошел
		for k := range occupiedTunnels {
			delete(occupiedTunnels, k)
		}

		// Проходим по каждому муравью
		for i := range ants {
			ant := &ants[i]

			// Если муравей уже на конечной позиции, пропускаем его
			if ant.Position >= len(ant.Path)-1 {
				continue
			}

			// Проверяем, можно ли переместить муравья
			nextPosition := ant.Position + 1
			currentRoom := ant.Path[ant.Position]
			nextRoom := ant.Path[nextPosition]
			tunnel := currentRoom + "-" + nextRoom

			if !occupiedRooms[nextRoom] && !occupiedTunnels[tunnel] {
				// Двигаем муравья, если следующая комната свободна и туннель не занят
				ant.Position = nextPosition
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))

				// Если муравей не входит в конечную комнату, отмечаем комнату и туннель как занятые
				if nextRoom != endRoom {
					occupiedRooms[nextRoom] = true
				}
				occupiedTunnels[tunnel] = true
				allAntsFinished = false
			}
		}

		if len(moves) > 0 {
			fmt.Printf("Turn %d: %s\n", turn+1, strings.Join(moves, " "))

			// fmt.Printf("%s\n", strings.Join(moves, " "))
		}

		turn++
		if turn > 100000 { // Защита от бесконечного цикла
			fmt.Println("Breaking to prevent a potential infinite loop.")
			break
		}

		// Сбрасываем занятость комнат, кроме стартовой и конечной
		for room := range occupiedRooms {
			if room != startRoom && room != endRoom {
				delete(occupiedRooms, room)
			}
		}
	}

	fmt.Println("Simulation ends")
}
