
<style type="text/css">

body{
    font-family      : helvetica, arial, freesans, clean, sans-serif;
    color            : #003269;
    background-color : #fff;
    border-color     : #999999;
    border-width     : 2px;
    line-height      : 1.5;
    margin           : 2em 3em;
    text-align       :left;
    font-size        : 16px;
    padding          : 0 100px 0 100px;

    width         : 1024px;
    margin-top    : 0px;
    margin-bottom : 2em;
    margin-left   : auto;
    margin-right  : auto;
}

h1 {
    font-family : 'Gill Sans Bold', 'Optima Bold', Arial, sans-serif;
    color       : #577AD3;
    font-weight : 400;
    font-size   : 48px;
}
h2{
    margin-bottom : 1em;
    padding-top   : 0.5em;
    color         : #003269;
    font-size     : 36px;
}
h3{
    border-bottom : 1px dotted #aaa;
    color         : #4660A4;
    font-size     : 30px;
}
h4 {
    font-size: 24px;
}
h5 {
    font-size: 18px;
}
code {
    font-family      : Consolas, "Inconsolata", Menlo, Monaco, Lucida Console, Liberation Mono, DejaVu Sans Mono, Bitstream Vera Sans Mono, Courier New, monospace, serif; /* Taken from the stackOverflow CSS*/
    background-color : #f5f5f5;
    border           : 1px solid #e1e1e8;
}


pre {
    display          : block;
    background-color : #f5f5f5;
    border           : 1px solid #ccc;
    padding          : 3px 3px 3px 3px;
}
pre code {
    white-space      : pre-wrap;
    padding          : 0;
    border           : 0;
    background-color : code;
}

table {
	border-collapse: collapse; border-spacing: 0;
	width: 100%;
	margin-bottom : 3em;
}
td, th {
	vertical-align: top;
	padding: 4px 10px;
	border: 1px solid #9BC3EB;
}
tr:nth-child(even) td, tr:nth-child(even) th {
	background: #EBF4FE;
}
th:nth-child(4) {
	width: auto;
}

</style>
# rello

## rello.proto

### Messages

<a name="Empty"></a>

#### Empty


<a name="ChecklistUpdate"></a>

#### ChecklistUpdate

Response structure from trello webhook
 TODO: Think about watching card in order to get card description and card
 due date for occurrence.Data and occurrence.Datetime respectively.
 TODO: Add references to the trello api documentation for fields

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| model | [Model](#Model) | 1 |  |
| action | [Action](#Action) | 2 |  |

<a name="Model"></a>

#### Model

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| id | TYPE_STRING | 1 |  |
| name | TYPE_STRING | 2 |  |
| idBoard | TYPE_STRING | 3 |  |
| idCard | TYPE_STRING | 4 |  |
| pos | TYPE_INT32 | 5 |  |
| checkItems | [CheckItem](#CheckItem) | 6 |  |

<a name="CheckItem"></a>

#### CheckItem

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| state | TYPE_STRING | 1 |  |
| idChecklist | TYPE_STRING | 2 |  |
| id | TYPE_STRING | 3 |  |
| name | TYPE_STRING | 4 |  |
| pos | TYPE_INT32 | 6 |  |

<a name="Action"></a>

#### Action

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| id | TYPE_STRING | 1 |  |
| idMemberCreator | TYPE_STRING | 2 |  |
| data | [Data](#Data) | 3 |  |
| type | TYPE_STRING | 4 |  |
| date | TYPE_STRING | 5 |  |
| memberCreator | [MemberCreator](#MemberCreator) | 6 |  |

<a name="Data"></a>

#### Data

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| checkItem | [CheckItem](#CheckItem) | 1 |  |
| checklist | [CheckList](#CheckList) | 2 |  |
| card | [Card](#Card) | 3 |  |
| board | [Board](#Board) | 4 |  |

<a name="CheckList"></a>

#### CheckList

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| name | TYPE_STRING | 1 |  |
| id | TYPE_STRING | 2 |  |

<a name="Card"></a>

#### Card

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| name | TYPE_STRING | 1 |  |
| id | TYPE_STRING | 2 |  |
| shortLink | TYPE_STRING | 3 |  |
| idShort | TYPE_INT32 | 4 |  |

<a name="Board"></a>

#### Board

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| name | TYPE_STRING | 1 |  |
| id | TYPE_STRING | 2 |  |
| shortLink | TYPE_STRING | 3 |  |

<a name="MemberCreator"></a>

#### MemberCreator

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| id | TYPE_STRING | 1 |  |
| avatarHas | TYPE_STRING | 2 |  |
| fullName | TYPE_STRING | 3 |  |
| initials | TYPE_STRING | 4 |  |
| username | TYPE_STRING | 5 |  |

### Services

#### Rello

Rello will need an application secret key, in order to make sure
 incoming requests are authenticated. Middleware!?

| Method Name | Request Type | Response Type | Description|
| ---- | ---- | ------------ | -----------|
| CheckListWebhook | ChecklistUpdate | Empty | CheckListWebhook is an endpoint for trello checklist webhooks that will
 parse the incoming data. It ensures that a request is from trello by doing
 a base64 encoding with the users secrete key.
 it then calls out to ambition-model to create an occurrence |

#### Rello - Http Methods

##### POST `/`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |
| model | body | [Model](#Model) |
| action | body | [Action](#Action) |

