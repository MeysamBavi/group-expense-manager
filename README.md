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
+ The structure of tables, including all headers are fixed. The *debt matrix* and *settlements* will be overwritten without being read so your changes are ignored.  
+ You can add or remove any expense or transactions from theirs sheets. But **be careful**: Removing an expense or transaction means **it never happened**.  
+ You can also edit the card number for each member but **don't edit their name**; Because the previous expenses and transactions use their previous name and cause error when loaded by the program.  
+ You can edit the *base state* matrix anytime but **remember**: *Base state* is the state of debts **before** running the create command and naturally **does not need editing**.

## Under the hood
*GEM* uses [excelize](https://github.com/qax-os/excelize) to create and edit the spreadsheets.
