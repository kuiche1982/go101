https://scene-si.org/2018/02/28/interfaces-in-go/


Interfaces in Go
Interfaces in Go are a powerful language construct, which allow you to define and use objects of different types under a common umbrella interface. What does that even mean? Imagine a Car, or a Motorbike. Both have a license plate. When we declare an interface with a License() string function, both Car and Motorbike objects will satisfy this interface.

package main

import (
    "fmt"
)

type Car struct {
    license string
}

func (c *Car) Name() string {
    return "car"
}
func (c *Car) License() string {
    return c.license
}

type MotorBike struct {
    license string
}

func (mb *MotorBike) Name() string {
    return "motor bike"
}
func (mb *MotorBike) License() string {
    return mb.license
}

type Vehicle interface {
    License() string
    Name() string
}

func PrintLicense(v Vehicle) {
    fmt.Println("I've seen a " + v.Name() + " with the license plate " + v.License())
}

func main() {
    car := Car{"LJ178FU"}
    bike := MotorBike{"LK6IDVR"}

    PrintLicense(&car)
    PrintLicense(&bike)
}
If you want to run the example, here’s the playground link for it. The output is what you might have expected however:

I've seen a car with the license plate LJ178FU
I've seen a motor bike with the license plate LK6IDVR
So far, you’re probablly thinking that you haven’t really heard anything new, or even that the example is useless. You’d be right. While the example does demonstrate a possible use case, it doesn’t really shine a light onto what makes interfaces useful.

What is the problem with interfaces?
If you’ll note, the concrete object receiver functions are pointers. This means that it’s mandatory to pass the Car or MotorBike objects as pointers to PrintLicense;

prog.go:42:14: cannot use car (type Car) as type Vehicle in argument to PrintLicense:
    Car does not implement Vehicle (License method has pointer receiver)
But the parameter decleration for the argument is v Vehicle and not v *Vehicle. You’d expect you can pass car and use the *Vehicle as the receiver - but it isn’t really so. Vehicle is an interface which can hold a value, by using *Vehicle, you’re declaring a pointer to an interface.

prog.go:42:15: cannot use &car (type *Car) as type *Vehicle in argument to PrintLicense:
    *Vehicle is pointer to interface, not interface
You might be tripping yourself up as to what interfaces here actually are. If you ever wrote PHP, you would have an interface class, which is very similar to Go’s interface decleration:

// Declare the interface 'iTemplate'
interface iTemplate
{
    public function setVariable($name, $var);
    public function getHtml($template);
}
The difference between PHP and Go however is that in PHP you enforce interface implementation to the class by declaring it like class Template implements iTemplate. In Go, you can use any object as an interface value, as long as it can satisfy the requirements.

The implications are:

You don’t need to implement specific interfaces on your structs as the part of the declaration, and your objects will satisfy multiple declarations you don’t even know about. One the most simple and ubiquitous interfaces that exists in go is the Stringer interface from the fmt package. Here’s an example from Go Tour.
While an interface in other languages is a set of functions that an implementation must adhere to, in Go it is a concrete value that holds the object that satisfies the interface. Unfortunately this also means that this interface blurs the lines between a pointer or concrete value which it may hold. If you use func(Car) Name() string however, don’t worry - Go will take care of the indirection in this case (playground link).
You may still be exposing yourself to runtime errors however, if you find yourself in a scenario where you need to use casting. You should be fully aware that this will expose you to runtime errors where your objects might not satisfy some interface. If Go supported some semblance of implements keywords on structs, some level of compile-time safety could be achieved, at least for the intended use case.
A real world example of interface use?
Ok. Let me postulate that you might have several different crypto currency tokens in your possesion. The balance of these wallets may be available on several, different, blockchain explorer websites. For example, if you have some Ethereum (ETH) tokens, your wallet details will be on etherchain, and if you have some Vertcoin (VTC) tokens, your wallet will show up on the vertcoin block explorer. There are several websites with APIs that will give out your wallet balance. You might even be on a crypto currency exchange like Bittrex which will give you several wallet balances based on the tokens that you hold there.

If to limit ourself just to the available public APIs, these two websites for ETH and VTC return significantly different results when it comes to retrieving wallet data. In fact, the vertcoin block explorer doesn’t even return JSON:

Vertcoin example API wallet balance,
Etherscan example API wallet balance
Bittrex /account/getbalances example
Abstracting these into an interface however, is simple enough:

type Balance struct {
        Currency string
        Balance  decimal.Decimal
}

type WalletBalance interface {
        Balance() ([]Balance, error)
        Name() string
}
The WalletBalance interface must provide a Name() function that would return a string of the used provider (Bittrex, Vertcoin, Etherscan,…). For each of these, there should be a list of token balances returned. This interface supports exchanges which will hold more than one wallet balance.

When it comes to the actual implementation, the main impact is that main loop where these interfaces would be handled. As you would just provide an object that satisfies the WalletBalance interface, a lot of boilerplate code gets taken out as the result.

for _, wallet := range owner.Wallets {
        var provider WalletBalance
        var err error

        switch wallet.Type {
        case "bittrex":
                provider, err = Bittrex{}.New(wallet.Config)
        case "etherscan":
                provider, err = Etherscan{}.New(wallet.Config)
        case "vertcoin":
                provider, err = Vertcoin{}.New(wallet.Config)
        default:
                err = fmt.Errorf("Unknown wallet type: %s", wallet.Type)
        }
        if err != nil {
                return result, err
        }
        balances, err := provider.Balance()
        // here comes the code that does something with balances...
This is the example of the main loop, how an implementation could look like. Before refactoring for interfaces, individual cases in the switch statement were significantly duplicated. Depending on what you want to do with []Balance, it’s going to be a lot of code that’s almost the same between individual cases, which was now moved after the switch.

A lot of differences between concrete implementations get abstracted away. For example, the ETH balance is expressed as a fixed point integer. To get the actual ETH balance, you should divide this number:

divisor, _ := decimal.NewFromString("1000000000000000000")
balance, _ := decimal.NewFromString(resp.Result)
balance = balance.Div(divisor)
minimumBalance, _ := decimal.NewFromString("0.002")
if balance.LessThan(minimumBalance) {
        return []Balance{}, nil
}
return []Balance{{"ETH", balance}}, nil
This is something that is fully abstracted into the Balance() call from the implemented interface. The implementation itself makes sure that the returned values follow a specific structure (ie, Balance). Without interfaces, boilerplate code is written to produce the same result for each type.

In another example, before optimizing it away with interfaces, this was the code to fetch different cryptocurrency mining pool balances:

func (p *Pool) Report(ms *markets.Markets) string {
        p.Currency = strings.ToUpper(p.Currency)
        price, err := ms.Price(p.Currency)
        if err != nil {
                return err.Error()
        }
        if strings.Contains(p.Link, ".ethermine.org") {
                provider, err := Ethermine{}.New(p.Currency, p.Link)
                if err != nil {
                        return err.Error()
                }
                return provider.Report(price)
        }
        if strings.Contains(p.Link, ".suprnova.cc") {
                provider, err := Suprnova{}.New(p.Currency, p.Link)
                if err != nil {
                        return err.Error()
                }
                return provider.Report(price)
        }
        return ""
}
As you can see, there’s opportunity to de-duplicate some code with the use of interfaces. While the effect with the example seems negligible, you should be aware that if you implement logic for more mining pools, the longer this function will be. If you change the signature of Report func, you’d have to refactor more code instead of just updating a few lines - the interface definition and the invocation of the function.

Wrapping it up
Interfaces are great for abstracting your functionality. They do however result in one or two rules that will make it easier to use:

You should never pass your objects without a pointer
You should never pass your interfaces with a pointer
The main reason is that, whenever you read something like Add(v Vector) you don’t really know if v is an interface or an object. It can be both, but relying on the above coding conventions, you can be reasonably sure, that Vector would be an interface, and that *Vector would be a pointer to the object. This makes it possible to read the function signature and know if you’re dealing with an interface or an object.

The logic behind the rules is as follows: The pointers to interfaces add a useless layer of indirection, that doesn’t give you any benefits. And very rarely do we want to pass a copy of an object to a function. Even if you would pass an object and expect a copy, you would have to think about any pointer or slice members - slices are also pointers, and they would be passed as-is, not copying the underlying data.

Of course, the lines get blurred even more when it comes to passing slices as arguments. The slices are already pointers, and having []Vector or []*Vector would only deal with memory overhead and cause additional indirection. If you have many edge cases where the rules should not apply, you might figure out a way how to namespace your interfaces into their own package. I mean, nothing says interface better than interfaces.Vector, but it might be a bit longer to type out.

By now you know why (in the true spirit of RFCs), I wrote SHOULD and not MUST in the rules above.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.