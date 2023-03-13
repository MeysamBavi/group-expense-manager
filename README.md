# GEM: Group Expense Manager
A CLI program in Go for managing group expenses in spreadsheets

## What does it solve?
When you are in group of friends, coworkers etc. and constantly lending and borrowing money by paying the group expenses, figuring out *who-owes-how-much-to-whom* can be cumbersome. *GEM* solves this problem by **providing an organized spreadsheet** to put everything at one place and make you free of any calculation.  

## What does it do?
In the spreadsheet provided by *GEM* you only need to enter the group **expenses and transactions** and the rest is handled; The debt between each two members is shown in a matrix and the **minimum transactions needed for settlement** are calculated.  
The initial state of your group **doesn't need to be even**. You can enter the current *base state* of the group; which is the current debt between each two members. This information will be used in later calculations.  

**Pro-tip**: upload the spreadsheet to a file sharing system like google sheets for everyone to have access!

## What do you mean by 'organized spreadsheet'?
In the very first moment, you give the names and card numbers of the group members to *GEM*, and it gives you back a spreadsheet, consisting of six **sheets**:

### Members
**Members** sheet contains the initial information you passed to program. Its main use is looking up someone's card number.

### Expenses
**Expenses** sheet contains the list of all expenses. You add a new row every time somebody pays for something.  
Each expense has a payer; The person who paid for the expense and lent money to the group. Each member has a *share weight* associated with that expense, showing how much of it is their share.  
*share weight* should be a non-negative integer or a boolean value (true is equivalent to 1, false to 0). A zero *share weight* means that member is not included in the expense.  
*Share Amount* is calculated via an Excel formula based on total amount, sum of *share weight*s and the member's *share weight*.

### Transactions
**Transactions** sheet contains the list of all transactions. To state that you have paid some of your debts to the group, add a new row.  
Each transaction has a *receiver*. The amount of transaction will be reduced from your overall debt and the debt state between you and *receiver* will be updated.

### Debt Matrix
**Debt Matrix** sheet contains the debt state between each two members. This matrix is calculated based on *expenses* *transactions* and *base state* **only** when you run the *update* command.  
For each cell, the person in the row should pay the person in the column. Only the positive values are shown in the matrix.

### Settlements
**Settlements** sheet the minimum transactions needed for settling up. This list is calculated based on *debt matrix* and **only** when you run the *update* command.

### Base State
**Base State** sheet contains the debt state between each two members, **before** creating the spreadsheet and using *GEM*. You can easily migrate to *GEM* by filling this matrix if you have been using a different system. The format of this matrix is similar to *debt matrix*.


## How do I use it?
Download the suitable binary for your system from [here](https://github.com/MeysamBavi/group-expense-manager/releases/latest). Create a spreadsheet by the **create** command:

```
gem create -o my-sheet-name.xlsx
```

Then you need to enter each member's name and card number. You can also pass the members' information in a `.csv` file like this:

```
gem create -o my-sheet-name.xlsx -f m.csv
```

And that's it. The spreadsheet is ready for entering the expenses and transactions.  
To add a new record, **copy and paste** the dummy row in *expenses*/*transactions* sheet to a new row. This will preserve the styling and Excel formulas!

After adding a few expenses or transactions, to calculate the debts, run the **update** command:

```
gem update my-sheet-name.xlsx --overwrite
```

This command will update the *debt matrix* and *settlement transactions* and overwrite the result on the same file. It calculates the debts based on **all** the expenses and transactions (and the *base state*).

This cycle is basically how you use *GEM*; Create the spreadsheet once, add some expenses and transactions, update the debts, add more expenses and transactions, update the debts again and so on.

You can run `gem help [command name]` to get all the details about a command, like its flags.

## What can I edit in the spreadsheet?

+ Generally the structure of tables, including all headers are fixed and not editable.
+ You can hide or unhide any sheets without any problem.

### Members
+ Values of *Card Number* column are editable.
+ The names are **not** editable; Because the old names will remain and be used all over the file.

### Expenses
+ Values of every column are editable except *Share Amount*.
+ *Share Amount* is calculated via an Excel formula, so don't edit it.
+ Members' names in the header are not editable.
+ **Be careful**; Removing an expense means **it never happened**.

### Transactions
+ Values of every column are editable.
+ **Be careful**; Removing a transaction means **it never happened**.

### Debt Matrix
+ *Debt Matrix* is **fully regenerated** with each *update* command and existing values are **ignored**.

### Settlements
+ *Settlements* are **fully regenerated** with each *update* command and existing values are **ignored**.

### Base State
+ Cell values are editable and read each time you run the *update* command.
+ Members' names in the margin are not editable.
+ **Reminder**: *Base state* is the state of debts **before** running the create command and naturally does not need editing.

## Under the hood
*GEM* uses [excelize](https://github.com/qax-os/excelize) to create and edit the spreadsheets.
