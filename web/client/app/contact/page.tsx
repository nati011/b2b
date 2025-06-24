'use client'

import { useState } from "react";
import { Mail, MapPin, PhoneIcon } from "lucide-react";

export default function Contact() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(false);
    const [success, setSuccess] = useState(false);
    const [formData, setFormData] = useState({
        name: "",
        email: "",
        subject: "",
        message: "",
    });

    async function sendEmail(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const apiEndpoint = "/api/email";
        try {
            setLoading(true);
            const response = await fetch(apiEndpoint, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(formData),
            });
        } catch (err) {
            setError(true);
            setLoading(false);
            setTimeout(() => {
                setError(false);
            }, 5000);
        }
    }

    return (
        <main className="mx-4">
            <div className="text-center lg:px-[8rem] grid grid-cols-1 gap-10 py-[6rem]">
                <div className=" text-center">
                    <p className="text-md font-semibold text-[#06516D]">Let&apos;s start a conversation. </p>
                    <p className="lg:text-[34px] text-lg font-bold text-[#032B3A]">
                        Get in touch with us.
                    </p>
                    <p className=" text-gray-600 my-4">
                        Have a project in mind or just want to learn more about how efoyeta store
                        can elevate your business? We&apos;re here to help. Reach out to us, and
                        let&apos;s start a conversation.
                    </p>
                </div>
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-[80px] my-10">
                    <div className="flex gap-4 items-center">
                        <div className="h-20 w-20 p-10 rounded-full bg-[#FFDD1F] flex justify-center items-center">
                            <span className="text-2xl text-[#043A4E]">
                                <Mail />
                            </span>
                        </div>
                        <div className="flex flex-col text-left ">
                            <p className="font-semibold text-[24px] text-[#043A4E]">Email us</p>
                            <p className="text-[18px] text-[#043A4E] font-light">
                                <a href="mailto:info@efoyeta store.com">
                                    info@efoyeta store.com
                                </a>
                            </p>
                        </div>
                    </div>

                    <div className="flex gap-4 items-center">
                        <div className="h-20 w-20 p-10 rounded-full bg-[#FFDD1F] flex justify-center items-center">
                            <span className="text-2xl text-[#043A4E]">
                                <PhoneIcon />
                            </span>
                        </div>
                        <div className="flex flex-col text-left">
                            <p className="font-semibold text-[24px] text-[#043A4E]">Call us</p>
                            <a className="text-[18px] font-light text-[#043A4E]" href="tel:+251912345678">+251912345678</a>
                        </div>
                    </div>

                    <div className="flex gap-4 items-center">
                        <div className="h-20 w-20 p-10 rounded-full bg-[#FFDD1F] flex justify-center items-center">
                            <span className="text-2xl text-[#043A4E]">
                                <MapPin />
                            </span>
                        </div>
                        <div className="flex flex-col text-left">
                            <p className="font-semibold text-[24px] text-[#043A4E]">Visit us</p>
                            <p className="text-[18px] font-light text-[#043A4E]">Addis Ababa, Ethiopia</p>
                        </div>
                    </div>

                </div>
                <div className="px-4 py-6 rounded-lg text-left bg-white border border-gray-100">
                    <form onSubmit={(e) => sendEmail(e)}>
                        <p className="text-[#032B3A] font-semibold">
                            Submit the form, and a member of our team will get back to you as
                            soon as possible.
                        </p>
                        <div className="my-4">
                            {error && <p className="text-red-500">An error occurred. Please try again.</p>}
                            {success && <p>Your message has been sent.</p>}
                        </div>
                        <div className="flex gap-2 flex-col items-left mt-4">
                            <div className="grid sm:grid-cols-3 gap-2">
                                <div className="">
                                    <label htmlFor="name" className="text-[#032B3A]">Name: </label>
                                    <input
                                        type="text"
                                        placeholder="John Doe"
                                        className="block w-full rounded-md border-0 py-1.5 pl-7 pr-20 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-0 sm:text-sm sm:leading-6  dark:bg-white dark:text-gray-900"
                                        id="name"

                                    />
                                </div>
                                <div className="w-full">
                                    <label htmlFor="email" className="text-[#032B3A]">Email: </label>
                                    <input
                                        type="email"
                                        placeholder="john.doe@example.com"
                                        className="block w-full rounded-md border-0 py-1.5 pl-7 pr-20 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-0 sm:text-sm sm:leading-6 dark:bg-white dark:text-gray-900"
                                        id="email"
                                        value={formData.email}
                                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                    />
                                </div>
                                <div className="">
                                    <label htmlFor="subject" className="text-[#032B3A]">Subject</label>
                                    <input
                                        type="text"
                                        placeholder="Subject"
                                        className="block w-full rounded-md border-0 py-1.5 pl-7 pr-20 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-0 sm:text-sm sm:leading-6  dark:bg-white dark:text-gray-900"
                                        id="subject"
                                        value={formData.subject}
                                        onChange={(e) => setFormData({ ...formData, subject: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="grid grid-cols-1">
                                <label htmlFor="message" className="text-[#032B3A]">Message: </label>
                                <textarea
                                    name=""
                                    id="message"
                                    cols={30}
                                    rows={10}
                                    className="p-2 border border-gray-300 focus:outline-none focus:ring-none rounded-lg  dark:bg-white dark:text-gray-900"
                                    placeholder="Write text here ..."
                                    value={formData.message}
                                    onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                                ></textarea>
                            </div>
                            <div className="grid grid-cols-1">
                                {
                                    loading ? (
                                        <input
                                            type="submit"
                                            value="Loading..."
                                            disabled={true}
                                            style={{ cursor: "not-allowed" }}
                                            className="p-2 border border-gray-300 focus:outline-none focus:ring-none rounded-lg bg-black text-white"
                                        />
                                    ) : <input
                                        type="submit"
                                        value="Submit"
                                        className="p-2 border border-gray-300 focus:outline-none focus:ring-none rounded-lg bg-black text-white"
                                    />
                                }
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </main>
    );
}
