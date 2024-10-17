import React from "react";
import { FaGoogle } from "react-icons/fa";

import { cn } from "~/lib/utils";

import { OttoLogo } from "~/components/branding/OttoLogo";
import { Button } from "~/components/ui/button";
import {
    Card,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "~/components/ui/card";

interface SignInProps {
    className?: string;
}

const SignIn: React.FC<SignInProps> = ({ className }) => {
    return (
        <Card className={cn("flex flex-col justify-between", className)}>
            <CardHeader>
                <CardTitle className="flex items-center justify-center">
                    <OttoLogo />
                </CardTitle>
                <CardDescription className="text-center w-3/4 mx-auto pt-4">
                    Please sign in using the button below.
                </CardDescription>
            </CardHeader>
            <CardFooter className="border-t pt-4">
                <Button
                    variant="secondary"
                    className="w-full"
                    onClick={() => {
                        window.location.href = "/oauth2/start?rd=/admin/";
                    }}
                >
                    <FaGoogle className="mr-2" />
                    Sign In with Google
                </Button>
            </CardFooter>
        </Card>
    );
};

export default SignIn;
