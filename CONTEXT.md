# Snook

This context owns how a snooker hall rents its tables by the hour, sells food and drink to the players at them, and settles up. It is the glossary for the Snook product (`snook-api`, `snook-web`); identity terms (User, Client, System, Role) come from UM and are not redefined here.

## Language

### Tables and play

**Table**:
One physical table in the hall, with a Table Type, an hourly rate and a status: available, in use, or under maintenance.
_Avoid_: Court, lane, slot

**Table Type**:
The class of a Table that sets its hourly rate and which Promotions apply to it.
_Avoid_: Category, tier

**Session**:
One continuous rental of a Table by a party of players, from opening to closing. Table time, Orders, Promotion and Payments all hang off the Session; the bill is the Session's bill.
_Avoid_: Game, booking, visit, ticket, table session (say Session)

**Open**:
Starting a Session on an available Table, which makes the Table in use.
_Avoid_: Check-in, start

**Pause**:
Stopping the clock on a Session without closing it. Paused minutes are not charged.
_Avoid_: Hold, break

**Move**:
Continuing a Session on a different Table; the clock and Orders follow the players.
_Avoid_: Transfer (the money word elsewhere), swap

**Close**:
Ending a Session: the clock stops, the bill is computed, Payments are taken or a Tab is opened, and the Table becomes available.
_Avoid_: Checkout, finish, end

**Table Charge**:
What the Session owes for time: the Table's hourly rate over the played minutes, less paused minutes.
_Avoid_: Rent, time fee

**Booking**:
A reservation of a Table for a named customer on a date and time range, which may be confirmed, used, or cancelled. A Booking is not a Session until the Table is opened.
_Avoid_: Reservation (say Booking), appointment

### Food and drink

**Menu Item**:
Something the hall sells to players at a Table, in a Menu Category, with a price and a stock count in some unit. It is low on stock when the count falls under the hall's threshold.
_Avoid_: Product, SKU, dish

**Menu Category**:
A named group of Menu Items in display order.
_Avoid_: Section, group

**Order**:
A quantity of one Menu Item added to a Session at the price of the moment, with any discount on that line. Orders reduce the Menu Item's stock.
_Avoid_: Line, cart, ticket

**Food Total**:
The sum of a Session's Orders.
_Avoid_: F&B, order total

### Money

**Promotion**:
A time-bounded offer that applies to Sessions on named Table Types: free hours after a number of played hours, a percentage off, or an amount off. One Promotion may be applied to a Session.
_Avoid_: Discount (the resulting amount), deal, campaign

**Bill**:
What a Session owes when it closes: Table Charge plus Food Total, less the Promotion's discount and any manual discount.
_Avoid_: Invoice, receipt (the printed form), total

**Payment**:
Money received against a Session's Bill by one method, such as cash or PromptPay, with a note. A Bill may be settled by several Payments.
_Avoid_: Transaction, settlement

**Tab**:
A Bill, or the unpaid part of one, that a named customer owes the hall after the Session closed, with an optional due date. It is outstanding until payments bring it to paid.
_Avoid_: Creditor (the code word), debt, IOU, credit

**Expense**:
Money the hall spent, in a category, on a date, with a description.
_Avoid_: Cost, outgoing

**Hall Settings**:
The hall's identity on receipts: name, address, phone, tax ID, receipt footer and PromptPay ID.
_Avoid_: Company, config, system settings
