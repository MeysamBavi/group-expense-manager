# GEM: Group Expense Manager
A CLI program in Go for managing group expenses in spreadsheets 

## What does it solve?
When you are in group of friends, coworkers etc. and constantly lending and borrowing money by paying the group expenses, figuring out *who-owes-how-much-to-whom* can be cumbersome. *GEM* solves this problem by **providing an organized spreadsheet** to put everything at one place and make you free of any calculation. 

## What does it do?
In the spreadsheet provided by *GEM* you only need to enter the group **expenses and transactions** and the rest is handled; The debt between each two members is shown in a matrix and the **minimum transactions needed for settlement** are calculated.  
The initial state of your group **doesn't need to be even**. You can enter the current *base state* of the group; which is the current debt between each two members. This information will be used in later calculations.

## How do I use it?
Download the suitable binary for your system from [here](https://github.com/MeysamBavi/group-expense-manager/releases/latest). Create a spreadsheet by the **create** command.

```
gem create -o my-sheet-name.xlsx
```

Then you need to enter each member's name and card number. You can also pass the members' information in a `.csv` file.

```
gem create -o my-sheet-name.xlsx -f m.csv
```

And that's it. Enter the group expenses and transactions and whenever you need, run the **update** command.

```
gem update my-sheet-name.xlsx --overwrite
```

This command will update the *debt matrix* and *settlement transactions*.

### What is editable?
+ Generally the structure of tables, including all headers are fixed and not editable.
+ You can hide or unhide any sheets without problem.

#### Members
+ Values of *Card Number* column are editable.
+ The names are **not** editable; Because the old names will remain and be used all over the file.

#### Expenses
+ Values of every column are editable except *Share Amount*.
+ *Share Amount* is calculated via an Excel formula, so don't edit it.
+ Members' names in the header are not editable.
+ **Be careful**; Removing an expense means **it never happened**.

#### Transactions
+ Values of every column are editable.
+ **Be careful**; Removing a transaction means **it never happened**.

#### Debt Matrix
+ *Debt Matrix* is **fully regenerated** with each *update* command and existing values are **ignored**.

#### Settlements
+ *Settlements* are **fully regenerated** with each *update* command and existing values are **ignored**.

#### Base State
+ Cell values are editable and read each time you run *update* command.
+ Members' names in the margin are not editable.
+ **Reminder**: *Base state* is the state of debts **before** running the create command and naturally does not need editing.

## Under the hood
*GEM* uses [excelize](https://github.com/qax-os/excelize) to create and edit the spreadsheets.
