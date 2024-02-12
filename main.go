package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Order struct {
	orderNum                              int
	date, orderName, orderCategory        string
	orderQuantity, orderCost, orderAmount int
}

func main() {
	records, err := readData("table.csv")
	if err != nil {
		log.Fatal(err)
	}

	orders := make([]Order, 0)
	for _, record := range records {
		orderNum, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		orderQuantity, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatal(err)
		}
		orderCost, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatal(err)
		}
		orderAmount, err := strconv.Atoi(record[6])
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, Order{
			orderNum:      orderNum,
			date:          record[1],
			orderName:     record[2],
			orderCategory: record[3],
			orderQuantity: orderQuantity,
			orderCost:     orderCost,
			orderAmount:   orderAmount,
		})
	}

	orders = quick_sort(orders, func(orders []Order, i, j int) bool {
		return orders[i].orderQuantity <= orders[j].orderQuantity
	})
	maxQuantityItem := orders[len(orders)-1].orderName

	totalRevenues := 0
	for _, order := range orders {
		totalRevenues += order.orderAmount
	}

	orders = quick_sort(orders, func(orders []Order, i, j int) bool {
		return orders[i].orderAmount <= orders[j].orderAmount
	})
	maxAmountItem := orders[len(orders)-1].orderName

	fmt.Printf("Общая выручка: %d, Товар, который был продан наибольшее количество раз: %s, Товар, который принес наибольшую выручку: %s\n",
		totalRevenues, maxQuantityItem, maxAmountItem)

	report := make([][]string, 0)
	report = append(report, []string{"Название товара", "Количество проданных единиц", "Доля в общей выручке"})

	for _, order := range orders {
		itemQuantity := strconv.Itoa(order.orderQuantity)
		totalRevenueShare := strconv.FormatFloat(float64(order.orderAmount)/float64(totalRevenues), 'f', 3, 32)
		report = append(report, []string{order.orderName, itemQuantity, totalRevenueShare})
	}

	f, err := os.Create("report.csv")
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	if err := w.WriteAll(report); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	totalRevenuesString := strconv.Itoa(totalRevenues)
	if err := w.Write([]string{"Общая выручка:", totalRevenuesString}); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}

func quick_sort(orders []Order, condition func(orders []Order, i, j int) bool) []Order {
	result := make([]Order, 0)
	less_arr := make([]Order, 0)
	greater_arr := make([]Order, 0)
	n := len(orders)
	if n == 0 {
		return result
	}
	if n == 1 {
		result = append(result, orders...)
		return result
	}
	pivot := n / 2
	for i := 0; i < n; i++ {
		if i != pivot {
			if condition(orders, i, pivot) {
				less_arr = append(less_arr, orders[i])
			} else {
				greater_arr = append(greater_arr, orders[i])
			}
		}
	}
	result = append(result, quick_sort(less_arr, condition)...)
	result = append(result, orders[pivot])
	result = append(result, quick_sort(greater_arr, condition)...)
	return result
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	r.Comma = ';'

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
