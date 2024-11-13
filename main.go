package grizzly

import "fmt"

func main() {
	numbers := []float64{10, 20, 30, 40, 50}
	names := []string{"David", "Roger", "Juliana", "Diana", "Sebastian"}
	floats := []string{"1.23", "4.56", "7.89", "10.12", "8.9"}
	test_serie := NewFloatSeries("test_1", numbers)
	test_series_2 := NewStringSeries("test_2", names)
	test_series_3 := NewStringSeries("test_3", floats)
	test_serie.Print(5)

	test_dataframe := CreateDataFrame(test_serie, test_series_2)
	test_dataframe.AddSeries(test_series_3)
	test_dataframe.Print(5)

	test_series_2.ReplaceWholeWord("Sebastian", "Sebastiano")
	test_series_2.Replace("D", "d")
	test_series_2.Print(5)

	filterFloat := func(value float64) bool {
		return value > 25
	}
	test_dataframe.FilterFloat("test_1", filterFloat)
	test_dataframe.Print(5)

	filterString := func(value string) bool {
		return value != "diana"
	}
	test_dataframe.FilterString("test_2", filterString)
	test_dataframe.Print(5)

	name := test_dataframe.GetColumnByName("test_3")
	name.ConvertStringToFloat()
	name.ConvertStringToFloat()

	empty := test_dataframe.CountEmpty()
	empty.Print(5)

	empty = test_dataframe.GetMax()
	empty.Print(5)

	empty = test_dataframe.GetMin()
	empty.Print(5)

	empty = test_dataframe.GetMean()
	empty.Print(5)

	empty = test_dataframe.GetMedian()
	empty.Print(5)

	empty = test_dataframe.GetVariance()
	empty.Print(5)

	elements := []string{"apple", "banana", "apple", "orange", "banana", "banana", "grape", "apple"}
	result := ArrayCountStringDuplicates(elements)
	fmt.Println("Duplicates:", result)
}
