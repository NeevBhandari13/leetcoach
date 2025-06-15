package problems

import "math/rand"

type Problem struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Difficulty string `json:"difficulty"`
}

var ProblemList []Problem

var twoSum = Problem{
	ID:         1,
	Title:      "Two Sum",
	Text:       "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
	Difficulty: "Easy",
}

var addTwoNumbers = Problem{
	ID:         2,
	Title:      "Add Two Numbers",
	Text:       "You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list. You may assume the two numbers do not contain any leading zero, except the number 0 itself.",
	Difficulty: "Medium",
}

// init function runs as soon as package is loaded up
func init() {
	ProblemList = append(ProblemList, twoSum, addTwoNumbers)
}

// function to get a random problem
func GetProblemText() string {
	randomIndex := rand.Intn(len(ProblemList))
	var problemText string = ProblemList[randomIndex].Text
	return problemText
}
